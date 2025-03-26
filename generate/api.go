package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/efucloud/common"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"

	"net/http"
	"os"
	"path"
	"reflect"
	"strings"
	"text/template"
	"time"
)

const apiHasPathParamsTemplate = "_{{_ .description _}}_" +
	"export const _{{_ .functionName _}}_  = async (params?: any) => {\n" +
	"    _{{_ if .extractPathParams _}}_ _{{_  .extractPathParams _}}_ _{{_ end _}}_\n" +
	"    return request_{{_ if .responseModel _}}_<API._{{_ .responseModel _}}_>_{{_ end _}}_(`_{{_ .api _}}_`, { method: '_{{_ .method _}}_' _{{_ if .params _}}_ , params: _{{_ .params _}}_ _{{_ end _}}_ _{{_ if .body _}}_ , body: _{{_ .body _}}_ _{{_ end _}}_ });\n" +
	"};\n"

const GlobalApiName = "GlobalApiName"
const ident = "    "

type RestAPI struct {
	// 从接口中自动提取
	structTypes        map[string]reflect.Type
	routes             []restful.Route
	apis               map[string]ApiData
	globalApiName      string
	generateTypescript bool
	files              map[string]string
	ignores            map[string]string
}

type ApiData struct {
	DocumentName  string
	RequestModel  string
	ResponseModel string
	Name          string
	Doc           string
	Notes         string
	Path          string
	Method        string
	Parameters    map[string]Parameters
	Response      map[int]string
}

func (api ApiData) String() string {
	return api.Method + " " + api.Path

}

type Parameters struct {
	Name        string // 参数名
	DataType    string
	Position    string // query path
	Description string
	Required    bool
	Enum        []string
	Default     string
}

func NewRestAPI(frontApiName string, generateTypescript bool) *RestAPI {
	api := &RestAPI{
		structTypes:        make(map[string]reflect.Type),
		apis:               make(map[string]ApiData),
		globalApiName:      frontApiName,
		generateTypescript: generateTypescript,
		files:              make(map[string]string),
		ignores:            make(map[string]string),
	}
	if len(api.globalApiName) == 0 {
		api.globalApiName = GlobalApiName
	}
	return api
}
func (rest *RestAPI) AddRoute(route restful.Route) {
	rest.routes = append(rest.routes, route)
}
func (rest *RestAPI) AddStruct(st reflect.Type) {
	rest.structTypes[st.Name()] = st
}
func (rest *RestAPI) AddStructIgnores(ignores ...string) {
	for _, item := range ignores {
		rest.ignores[item] = item
	}
}

func GetStructFieldDescription(item reflect.Type) string {
	result, _ := json.Marshal(ExtractStructFieldDescription(item))
	return string(result)
}
func ExtractStructFieldDescription(item reflect.Type) (result map[string]string) {
	result = make(map[string]string)
	for i := 0; i < item.NumField(); i++ {
		jsonName := item.Field(i).Tag.Get("json")
		if len(jsonName) == 0 {
			jsonName = item.Field(i).Name
		}
		description := item.Field(i).Tag.Get("description")
		if len(description) == 0 {
			description = item.Field(i).Name
		}
		result[jsonName] = description
	}
	return
}

const request = "import { request } from '@umijs/max';\n"

func (rest *RestAPI) GenerateToDir(dir string) {
	typeContent := rest.Generate()
	for k, c := range rest.files {
		_ = os.WriteFile(path.Join(dir, "api."+k+".ts"), []byte(request+c), os.ModePerm)
	}
	if len(typeContent) > 0 {
		_ = os.WriteFile(path.Join(dir, "types.d.ts"), []byte(typeContent), os.ModePerm)
	}
}
func (rest *RestAPI) Generate() (typeContent string) {

	rest.ParserRoutes()
	// 生成typescript定义
	if rest.generateTypescript {
		typs := NewTypeScript()
		for _, item := range rest.structTypes {
			typs.AddStruct(item)
		}
		for _, item := range rest.ignores {
			typs.AddStructIgnores(item)
		}
		typeContent = typs.Generate()
	}
	// 生成api
	for _, api := range rest.apis {
		rest.files[api.DocumentName] += rest.generateOneApi(api)
	}

	return
}
func (rest *RestAPI) generateOneApi(api ApiData) (content string) {
	var description string

	if len(api.Doc) > 0 {
		description += fmt.Sprintf("// %s\n", api.Doc)
	}
	if api.Doc != api.Notes && len(api.Notes) > 0 {
		description += fmt.Sprintf("// %s\n", api.Notes)
	}
	description += fmt.Sprintf("// 请求方法: %s\n", api.Method)
	description += fmt.Sprintf("// 请求地址: %s\n", api.Path)
	for code, data := range api.Response {
		description += fmt.Sprintf("// 响应码: %d  响应数据: %s\n", code, data)
	}
	var pathParams []string
	for name, item := range api.Parameters {
		required := "否"
		if item.Required {
			required = "是"
		}
		other := ""
		if len(item.Default) > 0 {
			other += fmt.Sprintf("默认值: %s", item.Default)
		}
		if len(item.Enum) > 0 {
			other += fmt.Sprintf(" 可选值: %s", strings.Join(item.Enum, ";"))
		}
		description += fmt.Sprintf("// 参数名: %s 参数类型: %s 参数位置: %s 是否必须: %s 参数说明: %s %s\n", name, item.DataType, item.Position, required, item.Description, other)
		if item.Position == "path" {
			pathParams = append(pathParams, item.Name)
		}
	}
	params := make(map[string]interface{})
	params["description"] = description
	params["functionName"] = api.Name
	params["method"] = strings.ToLower(api.Method)
	params["api"] = api.Path
	if len(pathParams) > 0 {
		params["extractPathParams"] = fmt.Sprintf("const { %s, ...rest } = params;", strings.Join(pathParams, ", "))
	}
	for _, p := range pathParams {
		params["api"] = strings.ReplaceAll(params["api"].(string), fmt.Sprintf("{%s}", p), fmt.Sprintf("${%s}", p))
	}
	switch api.Method {
	case http.MethodPost, http.MethodPut:
		if len(pathParams) > 0 {
			params["body"] = "rest"
		} else {
			params["body"] = "params"
		}
	case http.MethodGet, http.MethodDelete:
		if len(pathParams) > 0 {
			params["params"] = "rest"
		} else {
			params["params"] = "params"
		}
	}
	if len(api.ResponseModel) > 0 && !common.StringKeyInArray(api.ResponseModel, []string{"string", "uint", "bool", "float64"}) {
		params["responseModel"] = api.ResponseModel
	}
	t, _ := template.New(time.Now().String()).Delims("_{{_", "_}}_").Parse(apiHasPathParamsTemplate)
	b := new(bytes.Buffer)
	err := t.Execute(b, params)
	if err == nil {
		content += b.String()
	}
	return content
}
func (rest *RestAPI) ParserRoutes() {
	for _, route := range rest.routes {
		var api ApiData
		api.Parameters = make(map[string]Parameters)
		api.Response = make(map[int]string)
		api.Doc = route.Doc
		api.Notes = route.Notes
		api.Path = route.Path
		api.Method = route.Method
		if name, exist := route.Metadata[rest.globalApiName]; exist {
			api.Name = fmt.Sprintf("%v", name)
			api.Name = strings.ReplaceAll(api.Name, "[", "")
			api.Name = strings.ReplaceAll(api.Name, "]", "")
		} else {
			// todo 驼峰
			api.Name = fmt.Sprintf("%s%s", strings.ToLower(route.Method), route.Operation)
		}
		if doc, ex := route.Metadata[restfulspec.KeyOpenAPITags]; ex {
			n := strings.ReplaceAll(fmt.Sprintf("%v", doc), "[", "")
			n = strings.ReplaceAll(n, "]", "")
			api.DocumentName = strings.ReplaceAll(n, "-", "_")
		} else {
			api.DocumentName = "api"
		}
		for _, param := range route.ParameterDocs {
			var p Parameters
			p.Description = param.Data().Description
			p.Name = param.Data().Name
			p.DataType = param.Data().DataType
			p.Required = param.Data().Required
			p.Default = param.Data().DefaultValue
			p.Enum = param.Data().PossibleValues
			switch param.Kind() {
			case restful.PathParameterKind:
				p.Position = "path"
			case restful.QueryParameterKind:
				p.Position = "query"
			case restful.BodyParameterKind:
				p.Position = "body"
				if strings.Contains(p.DataType, ".") {
					dt := strings.Split(p.DataType, ".")
					if len(dt) >= 2 {
						p.DataType = dt[len(dt)-1]
					}
				}
			case restful.HeaderParameterKind:
				p.Position = "header"
			case restful.FormParameterKind:
				p.Position = "form"
			case restful.MultiPartFormParameterKind:
				p.Position = "multipart/form-data"
			}
			api.Parameters[p.Name] = p
			if route.ReadSample != nil {
				read := reflect.ValueOf(route.ReadSample).Type()
				if read != nil && read.Kind() != reflect.String {
					if read.Kind() != reflect.Slice {
						k := strings.TrimPrefix(read.Name(), "*")
						if strings.Contains(k, ".") {
							sp := strings.Split(k, ".")
							spLen := len(sp)
							k = sp[spLen-1]
						}
						rest.structTypes[read.Name()] = read
						api.RequestModel = k
					}
				}
			}
			if route.WriteSample != nil {
				write := reflect.ValueOf(route.WriteSample).Type()
				if write != nil && write.Kind() != reflect.String {
					if write.Kind() != reflect.Slice {
						k := strings.TrimPrefix(write.Name(), "*")
						if strings.Contains(k, ".") {
							sp := strings.Split(k, ".")
							spLen := len(sp)
							k = sp[spLen-1]
						}
						rest.structTypes[write.Name()] = write
						api.ResponseModel = k
					}
				}
			}
			for _, res := range route.ResponseErrors {
				if res.Code == http.StatusOK || res.Code == http.StatusCreated {
					if res.Model != nil {
						successRes := reflect.ValueOf(res.Model).Type()
						if len(api.ResponseModel) == 0 {
							api.ResponseModel = successRes.Name()
						}
						if successRes != nil {
							if successRes.Kind().String() == reflect.Pointer.String() || successRes.Kind().String() == reflect.Struct.String() {
								rest.structTypes[successRes.Name()] = successRes
								api.Response[res.Code] = successRes.Name()
							}
						}
					}
				} else {
					if res.Model != nil {
						mo := reflect.ValueOf(res.Model).Type()
						if mo != nil {
							if mo.Kind().String() == reflect.Pointer.String() || mo.Kind().String() == reflect.Struct.String() {
								api.Response[res.Code] = GetStructFieldDescription(mo)
							}
						}
					}

				}
			}
		}
		rest.apis[api.String()] = api
	}
}
