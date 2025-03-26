/*
Copyright 2022 The efucloud.com Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package common

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/net/http2"
	"k8s.io/klog/v2"
	"net"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"time"
)

type ApiInfo struct {
	Tag         string
	Description string
}

var (
	ApiInfos []ApiInfo
)

func RegisterApiInfo(info ApiInfo) {
	ApiInfos = append(ApiInfos, info)
}
func (a ApiInfo) Tags() []string {
	return []string{a.Tag}
}

func GetRequestPaginationInformation(req *restful.Request) (current int, pageSize int, order string) {
	current = String2Int(req.QueryParameter("current"), DefaultPage)
	pageSize = String2Int(req.QueryParameter("pageSize"), DefaultPageSize)
	order = req.QueryParameter("order")
	if len(order) == 0 {
		order = DefaultOrder
	}
	return current, pageSize, order
}

type QueryParam struct {
	WhereQuery string
	WhereArgs  []interface{}
}

func CreateHttpClient(useHttp2 bool, timeout time.Duration) (client *http.Client) {
	if useHttp2 {
		client = &http.Client{
			Transport: &http2.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: timeout,
		}
	} else {
		client = &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: timeout,
		}
	}
	return client
}
func HttpRequest(client *http.Client, method, address string, headers, cookies map[string]interface{}, queries url.Values, body []byte) (response *http.Response, err error) {
	req, err := http.NewRequest(method, address, bytes.NewReader(body))
	if err != nil {
		err = fmt.Errorf("create http request failed, method: %s, address: %s, err: %s", method, address, err.Error())
		klog.Error(err)
		return response, err
	}
	for k, v := range headers {
		req.Header.Add(k, fmt.Sprintf("%s", v))
	}
	for k, v := range cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: fmt.Sprintf("%s", v)})
	}
	req.URL.RawQuery = queries.Encode()
	return client.Do(req)
}

type ResponseError struct {
	Message    string        `json:"message" yaml:"message" description:"错误英文编码"`
	Detail     string        `json:"detail" yaml:"detail" description:"错误详情信息"`
	Links      []ErrorSource `json:"links,omitempty" description:""`
	Alert      string        `json:"alert" yaml:"alert" description:"支持I18N的提示信息"`
	RequestURI string        `json:"requestUri" description:"当前请求地址"`
}
type ErrorSource struct {
	File string `json:"file" description:""`
	Line int    `json:"line" description:""`
}
type AuthRedirectInfo struct {
	Message               string                 `json:"message" description:"提示信息"`
	Params                map[string]interface{} `json:"params" description:"重定向参数"`
	AuthorizationEndpoint string                 `json:"authorizationEndpoint" description:"认证提供商地址"`
	Alert                 string                 `json:"alert" description:"提示信息"`
}

type ResponseList struct {
	Data  any   `json:"data" yaml:"data"`
	Total int64 `json:"total" yaml:"total"`
}

func ResponseSuccess(resp *restful.Response, info interface{}) {
	resp.Header().Add("X-Content-Type-Options", "nosniff")
	resp.Header().Add("X-XSS-Protection", "1; mode=block")
	_ = resp.WriteAsJson(info)

}
func ResponseAuthRedirect(ctx context.Context, resp *restful.Response, bundle *i18n.Bundle, lang, message, authorizationEndpoint string,
	params map[string]interface{}) {
	var body AuthRedirectInfo
	body.Message = message
	body.Params = params
	body.AuthorizationEndpoint = authorizationEndpoint
	body.Alert, _ = GetLocaleMessage(bundle, nil, lang, "statusUnauthorized")
	_ = resp.WriteHeaderAndJson(http.StatusUnauthorized, body, restful.MIME_JSON)
}
func ResponseErrorMessage(ctx context.Context, req *restful.Request, resp *restful.Response, bundle *i18n.Bundle, detail ErrorData) {
	if detail.ResponseCode == 0 {
		detail.ResponseCode = http.StatusInternalServerError
	}

	resp.Header().Add("X-Content-Type-Options", "nosniff")
	resp.Header().Add("X-XSS-Protection", "1; mode=block")
	var body ResponseError
	body.Message = detail.MsgCode
	if detail.Err != nil {
		body.Detail = detail.Err.Error()
	}
	depth := 1

	for {
		if _, file, line, ok := runtime.Caller(depth); ok {
			body.Links = append(body.Links, ErrorSource{File: file, Line: line})
		} else {
			break
		}
		depth += 1
		if depth > 5 {
			break
		}
	}

	body.RequestURI = req.Request.RequestURI
	body.Alert, _ = GetLocaleMessage(bundle, detail.Params, detail.Lang, detail.MsgCode)
	_ = resp.WriteHeaderAndJson(detail.ResponseCode, body, restful.MIME_JSON)

}

// RequestQuery paramType: string,number queryType: eq,like
func RequestQuery(name, paramType, queryType string, req *restful.Request, queryParam *QueryParam) {
	if queryType == QueryTypeIn {
		valueSlice := req.QueryParameters(fmt.Sprintf("%s[]", name))
		if paramType == ParamTypeStringSlice {
			if queryParam.WhereQuery == "" {
				queryParam.WhereQuery = fmt.Sprintf(" %s IN (?) ", CamelString2Snake(name))
			} else {
				queryParam.WhereQuery += fmt.Sprintf(" AND %s IN (?)", CamelString2Snake(name))
			}
			queryParam.WhereArgs = append(queryParam.WhereArgs, valueSlice)
		} else if paramType == ParamTypeNumber {
			if queryParam.WhereQuery == "" {
				queryParam.WhereQuery = fmt.Sprintf(" %s IN (?) ", CamelString2Snake(name))
			} else {
				queryParam.WhereQuery += fmt.Sprintf(" AND %s IN (?)", CamelString2Snake(name))
			}
			queryParam.WhereArgs = append(queryParam.WhereArgs, StringsToUints(valueSlice))
		}
	} else {
		if strings.HasPrefix(name, "search:") {
			value := req.QueryParameter("search")
			nameList := strings.Split(strings.TrimPrefix(name, "search:"), ";")
			var sqlList []string
			for _, item := range nameList {
				sqlList = append(sqlList, fmt.Sprintf(" %s %s ?", CamelString2Snake(item), queryType))
				queryParam.WhereArgs = append(queryParam.WhereArgs, fmt.Sprintf("%%%s%%", value))
			}
			if queryParam.WhereQuery == "" {
				queryParam.WhereQuery = fmt.Sprintf("(%s)", strings.Join(sqlList, " OR "))
			} else {
				queryParam.WhereQuery += fmt.Sprintf(" AND ( %s )", strings.Join(sqlList, " OR "))
			}
		}
		value := req.QueryParameter(name)
		nv := strings.TrimSpace(value)
		if nv != "" {
			//相等可以是字符或者数字
			if queryType == "" || queryType == QueryTypeEqual {
				if paramType == ParamTypeNumber {
					v := StringsToUint(nv)
					if v > 0 {
						if queryParam.WhereQuery == "" {
							queryParam.WhereQuery = fmt.Sprintf(" %s = ? ", CamelString2Snake(name))
						} else {
							queryParam.WhereQuery += fmt.Sprintf(" AND %s = ? ", CamelString2Snake(name))
						}
						queryParam.WhereArgs = append(queryParam.WhereArgs, v)
					}
				} else if paramType == ParamTypeString || paramType == "" {
					if queryParam.WhereQuery == "" {
						queryParam.WhereQuery = fmt.Sprintf(" %s = ? ", CamelString2Snake(name))
					} else {
						queryParam.WhereQuery += fmt.Sprintf(" AND %s = ? ", CamelString2Snake(name))
					}
					queryParam.WhereArgs = append(queryParam.WhereArgs, nv)
				} else if paramType == ParamTypeBool {
					if queryParam.WhereQuery == "" {
						queryParam.WhereQuery = fmt.Sprintf(" %s = ? ", CamelString2Snake(name))
					} else {
						queryParam.WhereQuery += fmt.Sprintf(" AND %s = ? ", CamelString2Snake(name))
					}
					if strings.TrimSpace(nv) == "0" || strings.TrimSpace(strings.ToUpper(nv)) == "F" || strings.TrimSpace(strings.ToUpper(nv)) == "FALSE" {
						queryParam.WhereArgs = append(queryParam.WhereArgs, 0)
					} else {
						queryParam.WhereArgs = append(queryParam.WhereArgs, 1)
					}
				}
				// like 只能为字符串
			} else if queryType == QueryTypeLike {
				if queryParam.WhereQuery == "" {
					queryParam.WhereQuery = fmt.Sprintf(" %s LIKE  ? ", CamelString2Snake(name))
				} else {
					queryParam.WhereQuery += fmt.Sprintf(" AND %s LIKE  ? ", CamelString2Snake(name))
				}
				queryParam.WhereArgs = append(queryParam.WhereArgs, fmt.Sprintf("%%%s%%", nv))

			}
		}
	}
}

// RequestQuerySearch queryType: eq,like use 'or' to connect
func RequestQuerySearch(value, queryType string, fields []string, queryParam *QueryParam) {
	if len(value) == 0 || len(fields) == 0 {
		return
	}
	//相等可以是字符或者数字
	if queryType == "" || queryType == QueryTypeEqual {
		for _, name := range fields {
			if queryParam.WhereQuery == "" {
				queryParam.WhereQuery = fmt.Sprintf(" %s = ? ", CamelString2Snake(name))
			} else {
				queryParam.WhereQuery += fmt.Sprintf(" OR %s = ? ", CamelString2Snake(name))
			}
			queryParam.WhereArgs = append(queryParam.WhereArgs, value)
		}
		// like 只能为字符串
	} else if queryType == QueryTypeLike {
		for _, name := range fields {
			if queryParam.WhereQuery == "" {
				queryParam.WhereQuery = fmt.Sprintf(" %s LIKE  ? ", CamelString2Snake(name))
			} else {
				queryParam.WhereQuery += fmt.Sprintf(" OR %s LIKE  ? ", CamelString2Snake(name))
			}
			queryParam.WhereArgs = append(queryParam.WhereArgs, fmt.Sprintf("%%%s%%", value))
		}
	}
}

func QueryEqual(name string, value interface{}, queryParam *QueryParam) {
	if queryParam.WhereQuery == "" {
		queryParam.WhereQuery = fmt.Sprintf(" %s = ? ", CamelString2Snake(name))
	} else {
		queryParam.WhereQuery += fmt.Sprintf(" AND %s = ? ", CamelString2Snake(name))
	}
	queryParam.WhereArgs = append(queryParam.WhereArgs, value)
}

func QueryIn(name string, values interface{}, queryParam *QueryParam) {
	if queryParam.WhereQuery == "" {
		queryParam.WhereQuery = fmt.Sprintf(" %s IN (?) ", CamelString2Snake(name))
	} else {
		queryParam.WhereQuery += fmt.Sprintf(" AND %s IN (?) ", CamelString2Snake(name))
	}
	queryParam.WhereArgs = append(queryParam.WhereArgs, values)
}

func TokenErr(resp *restful.Response, typ, description string, statusCode int) error {
	data := struct {
		Error       string `json:"error"`
		Description string `json:"error_description,omitempty"`
	}{typ, description}
	resp.ResponseWriter.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(statusCode)
	_ = resp.WriteAsJson(data)
	return nil
}

func Request(method, address string, headers map[string]string, queries map[string]interface{}, body interface{}) (response *http.Response, err error) {
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second,
	}
	b := new(bytes.Buffer)
	err = json.NewEncoder(b).Encode(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, address, b)
	if err != nil {
		klog.Errorf("create request failed, err: %s", err.Error())
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	queryValues := url.Values{}
	for k, v := range queries {
		queryValues.Add(k, fmt.Sprintf("%v", v))
	}
	req.URL.RawQuery = queryValues.Encode()
	return client.Do(req)
}

func NewHTTPClientWithCA(rootCA string, insecureSkipVerify bool) (client *http.Client, err error) {

	var block *pem.Block
	block, _ = pem.Decode([]byte(rootCA))
	if block == nil {
		err = errors.New("ca decode failed")
		return nil, err
	}
	// Only use PEM "CERTIFICATE" blocks without extra headers
	if block.Type != "CERTIFICATE" || len(block.Headers) != 0 {
		err = fmt.Errorf("ca decode failed, block type: %s is not CERTIFICATE", block.Type)
		return nil, err

	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		err = fmt.Errorf("ca decode failed, err: %s", err.Error())
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AddCert(cert)
	// Copied from http.DefaultTransport.
	tlsConfig := tls.Config{RootCAs: pool, InsecureSkipVerify: insecureSkipVerify}
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}).DialContext,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
		},
	}
	return client, err
}

func GetLangFromCtx(ctx context.Context, key string) (lang string) {
	lang = "zh"
	if len(key) == 0 {
		key = "RequestLanguage"
	}
	lan := ctx.Value(key)
	if lan != nil {
		lang = lan.(string)
	}
	return lang
}
