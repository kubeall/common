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

type OpenIDConfiguration struct {
	Issuer                                 string   `json:"issuer" description:""`
	AuthorizationEndpoint                  string   `json:"authorization_endpoint" description:""`
	TokenEndpoint                          string   `json:"token_endpoint" description:""`
	UserinfoEndpoint                       string   `json:"userinfo_endpoint" description:""`
	JwksUri                                string   `json:"jwks_uri" description:""`
	ResponseTypesSupported                 []string `json:"response_types_supported" description:""`
	ResponseModesSupported                 []string `json:"response_modes_supported" description:""`
	GrantTypesSupported                    []string `json:"grant_types_supported" description:""`
	SubjectTypesSupported                  []string `json:"subject_types_supported" description:""`
	IdTokenSigningAlgValuesSupported       []string `json:"id_token_signing_alg_values_supported" description:""`
	ScopesSupported                        []string `json:"scopes_supported" description:""`
	ClaimsSupported                        []string `json:"claims_supported" description:""`
	RequestParameterSupported              bool     `json:"request_parameter_supported" description:""`
	RequestObjectSigningAlgValuesSupported []string `json:"request_object_signing_alg_values_supported" description:""`
}

func (OpenIDConfiguration) GormDataType() string {
	return "json"
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (ins *OpenIDConfiguration) Scan(value interface{}) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal OpenIDConfiguration value: ", value))
	}
	err := json.Unmarshal(byteValue, ins)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (ins OpenIDConfiguration) Value() (driver.Value, error) {
	re, err := json.Marshal(ins)
	return re, err
}
