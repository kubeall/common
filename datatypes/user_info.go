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

type UserInfoPath struct {
	Username     string `json:"username" yaml:"username" description:"用户名"`
	Email        string `json:"email" yaml:"email" description:"邮箱"`
	Nickname     string `json:"nickname" yaml:"nickname" description:"别名"`
	Organization string `json:"organization" yaml:"organization" description:"组织"`
	ID           string `json:"id" yaml:"id" description:"记录ID"`
}

func (UserInfoPath) GormDataType() string {
	return "json"
}

// Scan 实现 sql.Scanner 接口，Scan 将 value 扫描至 Jsonb
func (ins *UserInfoPath) Scan(value interface{}) error {
	byteValue, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal OpenIDConfiguration value: ", value))
	}
	err := json.Unmarshal(byteValue, ins)
	return err
}

// Value 实现 driver.Valuer 接口，Value 返回 json value
func (ins UserInfoPath) Value() (driver.Value, error) {
	re, err := json.Marshal(ins)
	return re, err
}
