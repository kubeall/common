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
	"crypto/tls"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"golang.org/x/net/http2"
	"k8s.io/klog/v2"
	"net/http"
	"net/url"
	"time"
)

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

func GetRequestPaginationInformation(req *restful.Request) (page int, size int, order string) {
	page = String2Int(req.QueryParameter("page"), DefaultPage)
	size = String2Int(req.QueryParameter("size"), DefaultPageSize)
	order = req.QueryParameter("order")
	if len(order) == 0 {
		order = DefaultOrder
	}
	return page, size, order
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
