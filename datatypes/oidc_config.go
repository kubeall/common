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

package datatypes

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
)

type OidcConfig struct {
	// 提供商的地址，如https://gitlab.com,在不配置Certificate的情况下，程序会根据https://gitlab.com/.well-known/openid-configuration获取token的公钥
	Issuer string `json:"issuer" yaml:"issuer" description:"提供商的地址"`
	// 应用的ClientID
	ClientID string `json:"clientId" yaml:"clientId" description:"应用的ClientID"`
	// 应用的ClientSecret
	ClientSecret string `json:"clientSecret" yaml:"clientSecret" description:"应用的ClientSecret"`
	// 跳转到认证的页面，如https://gitlab.com/oauth/authorize，该信息会返回给前端用于前端组成认证重定向地址
	AuthorizationEndpoint string `json:"authorizationEndpoint" yaml:"authorizationEndpoint" description:"跳转到认证的页面"`
	// 认证完成后的重定向地址，用于接收返回的code，如gitlab认证成功后返回的code,state或者err信息,
	// 前后端分离模式下，该地址为前端地址，可由前端自行拼接
	RedirectURI string `json:"redirectUri" yaml:"redirectUri" description:"认证完成后的重定向地址"`
	// 获取eauth Token的地址
	TokenEndpoint string `json:"tokenEndpoint" yaml:"tokenEndpoint" description:"Token获取地址"`
	// 获取用户信息的地址
	UserinfoEndPoint string `json:"userinfoEndPoint" yaml:"userinfoEndPoint" description:"用户信息的地址"`
	// 提供商的ca信息，可以不提供，
	IssuerCA      string `json:"issuerCa" yaml:"issuerCA" description:"提供商的ca信息"`
	UsernameClaim string `json:"usernameClaim" yaml:"usernameClaim" description:"用户名"`
	GroupsClaim   string `json:"groupsClaim" yaml:"groupsClaim" description:"组信息"`
	// token校验的公钥信息，若不配置，应用需要根据Issuer+/.well-known/openid-configuration去获取
	// 若以gitlab为例https://gitlab.com/.well-known/openid-configuration
	Certificate string   `json:"certificate" yaml:"certificate" description:"token校验的公钥信息"`
	Scopes      []string `json:"scopes" yaml:"scopes" description:"请求的域信息"`
}

func (OidcConfig) GormDataType() string {
	return "json"
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (ins *OidcConfig) Scan(value interface{}) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal OpenIDConfiguration value: ", value))
	}
	err := json.Unmarshal(byteValue, ins)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (ins OidcConfig) Value() (driver.Value, error) {
	re, err := json.Marshal(ins)
	return re, err
}
