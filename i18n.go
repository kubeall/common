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
	"embed"
	"github.com/ghodss/yaml"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"k8s.io/klog/v2"
	"path"
)

const (
	I18nZH = "zh"
	I18nEN = "en"
)

var (
	UniversalTranslator *ut.UniversalTranslator
	Bundle              *i18n.Bundle
)

func I18nInit(i18nFiles embed.FS) {
	// todo add other language
	UniversalTranslator = ut.New(en.New(), zh.New())
	Bundle = i18n.NewBundle(language.Chinese)
	Bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	_, err := Bundle.LoadMessageFileFS(i18nFiles, path.Join("locales", "en.yaml"))
	if err != nil {
		klog.Fatalf("load i18n message file failed, err: %s", err.Error())
	}
	_, err = Bundle.LoadMessageFileFS(i18nFiles, path.Join("locales", "zh.yaml"))
	if err != nil {
		klog.Fatalf("load i18n message file failed, err: %s", err.Error())
	}
}
func GetLocaleMessage(templateData map[string]interface{}, lang string, id string) (msg string, err error) {
	localizer := i18n.NewLocalizer(Bundle, lang)
	msg, err = localizer.Localize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: id}, TemplateData: templateData})
	if len(msg) == 0 {
		msg = id
	}
	return msg, err
}
