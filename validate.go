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
	"context"
	"fmt"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"regexp"
	"strings"
)

var dns1123Reg *regexp.Regexp

func init() {
	dns1123Reg = regexp.MustCompile(`[a-z0-9]([-a-z0-9]*[a-z0-9])?`)
}
func ValidateDNS1123(fl validator.FieldLevel) bool {
	return dns1123Reg.MatchString(fl.Field().String())
}
func NotBlank(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		return len(strings.TrimSpace(field.String())) > 0
	case reflect.Chan, reflect.Map, reflect.Slice, reflect.Array:
		return field.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !field.IsNil()
	default:
		return field.IsValid() && field.Interface() != reflect.Zero(field.Type()).Interface()
	}
}
func ValidateTrans(validate *validator.Validate, lang string, err error) FiledValidFailed {
	trans, _ := UniversalTranslator.GetTranslator(lang)
	switch lang {
	case I18nZH:
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)
	case I18nEN:
		_ = en_translations.RegisterDefaultTranslations(validate, trans)
	default:
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)
	}
	errs := err.(validator.ValidationErrors)
	return removeTopStruct(errs.Translate(trans))
}

func ValidateTransCtx(lang string, validate *validator.Validate, ctx context.Context, err error) FiledValidFailed {

	trans, _ := UniversalTranslator.GetTranslator(lang)
	switch lang {
	case I18nZH:
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)
	case I18nEN:
		_ = en_translations.RegisterDefaultTranslations(validate, trans)
	default:
		_ = zh_translations.RegisterDefaultTranslations(validate, trans)
	}
	errs := err.(validator.ValidationErrors)

	return removeTopStruct(errs.Translate(trans))
}
func TagNameFunc(fld reflect.StructField) string {
	name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
	if name == "-" {
		return fld.Name
	}
	return name
}

func removeTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}

type FiledValidFailed map[string]string

func (f FiledValidFailed) String() string {
	var infos []string
	for k, v := range f {
		infos = append(infos, fmt.Sprintf("%s:%s", k, v))
	}
	return strings.Join(infos, ";")
}
func (f FiledValidFailed) LocaleString(localeMap map[string]interface{}) string {
	var infos []string
	for key, value := range f {
		if v, exist := localeMap[key]; exist {
			infos = append(infos, fmt.Sprintf("%s:%s", v, value))
		} else {
			infos = append(infos, fmt.Sprintf("%s:%s", key, value))
		}
	}
	return strings.Join(infos, ";")
}
func (f FiledValidFailed) LocaleMap(localeMap map[string]string) (result map[string]string) {
	result = make(map[string]string)
	for key, value := range f {
		if v, exist := localeMap[key]; exist {
			result[v] = value
		} else {
			result[key] = value
		}
	}
	return result
}
