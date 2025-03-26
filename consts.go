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

const DefaultOrder = "id desc"
const DefaultPage = 1
const DefaultPageSize = 20
const (
	QueryTypeEqual       = "eq"
	QueryTypeLike        = "like"
	QueryTypeIn          = "in"
	ParamTypeString      = "string"
	ParamTypeNumber      = "integer"
	ParamTypeBool        = "bool"
	ParamTypeStringSlice = "stringSlice"
	ParamTypeNumberSlice = "numberSlice"
)
const (
	I18nZH = "zh"
	I18nEN = "en"
)
const TimeFormat = "2006-01-02 15:04:05"
const Enable = 1
const Disable = 0
