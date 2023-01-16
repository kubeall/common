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
	"crypto/md5"
	"crypto/rand"
	"encoding/base32"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"github.com/ghodss/yaml"
	"io"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
	"path"
	"reflect"
	"strconv"
	"strings"
)

func GetMachineID(app string) string {
	id, _ := machineid.ProtectedID(app)
	klog.Infof("generate app: %s machine id: %s", app, id)
	return id
}

func GetAllFiles(dirPath string) (dirs []string, files []string, err error) {
	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}
	for _, f := range fs {
		if f.IsDir() {
			dirs = append(dirs, path.Join(dirPath, f.Name()))
			ds, fs, er := GetAllFiles(path.Join(dirPath, f.Name()))
			if er == nil {
				dirs = append(dirs, ds...)
				files = append(files, fs...)
			}
		} else {
			files = append(files, path.Join(dirPath, f.Name()))
		}
	}

	return
}
func GetModuleFiles(dirPath string) (dirs []string, modules map[string][]string, err error) {
	modules = make(map[string][]string)
	modules[""] = []string{}
	fs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return
	}
	for _, f := range fs {
		if f.IsDir() {
			p := path.Join(dirPath, f.Name())
			modules[p] = []string{}
			dirs = append(dirs, p)
			ds, fs, er := GetAllFiles(path.Join(dirPath, f.Name()))
			if er == nil {
				dirs = append(dirs, ds...)
				modules[p] = fs
			}
		} else {
			modules[""] = append(modules[""], path.Join(dirPath, f.Name()))
		}
	}
	return
}

func LoadConfig(path string, object interface{}) {
	if strings.HasSuffix(path, ".yaml") || strings.HasSuffix(path, ".yml") {
		if data, err := ioutil.ReadFile(path); err == nil {
			if jsonData, err := yaml.YAMLToJSON(data); err != nil {
				klog.Fatalf("Unable to decode application %s yaml config from file, err: %s", path, err)
			} else {
				if err = json.Unmarshal(jsonData, object); err != nil {
					klog.Fatalf("Unable to decode application %s json config from file, err: %s", path, err)
				}
			}
		} else {
			klog.Fatalf("Unable to read application %s yaml config from file, err: %s", path, err)

		}
	} else if strings.HasSuffix(path, ".json") {
		if data, err := ioutil.ReadFile(path); err == nil {
			if err = json.Unmarshal(data, object); err != nil {
				klog.Fatalf("Unable to decode application %s json config from file, err: %s", path, err)
			}
		} else {
			klog.Fatalf("Unable to read application %s json config from file, err: %s", path, err)

		}
	} else {
		klog.Fatalf("Unable to read application  config from file: %s", path)
	}
}

func MD5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func KeyInArray(key string, arrays []string) (exist bool) {
	for _, v := range arrays {
		if key == v {
			exist = true
			return
		}
	}
	return
}
func String2Int(str string, defVal int) int {
	if in, err := strconv.Atoi(str); err != nil {
		return defVal
	} else {
		if in < 1 {
			in = 1
		}
		return in
	}

}
func StringsToUint(str string) (int uint) {

	if i, e := strconv.Atoi(str); e == nil {
		return uint(i)
	}
	return 0

}

// Snake2CamelString snake to camel string, xx_yy to XxYy
func Snake2CamelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true

			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
func CamelString2Snake(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')

		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}
func StringsToUints(strings []string) (ints []uint) {
	for _, str := range strings {
		if i, e := strconv.Atoi(str); e == nil {
			if i > 0 {
				ints = append(ints, uint(i))
			}

		}
	}
	return

}

func GetStructFieldsType(v interface{}) (fields map[string]string) {
	fields = make(map[string]string)
	dataType := reflect.TypeOf(v)
	dataValue := reflect.ValueOf(v)
	if dataType.Kind() == reflect.Ptr {
		originType := reflect.ValueOf(v).Elem().Type()
		if originType.Kind() != reflect.Struct {
			return
		}
		// 解引用
		dataType = dataType.Elem()
		dataValue = dataValue.Elem()
		num := dataType.NumField()
		for i := 0; i < num; i++ {
			name := strings.SplitN(dataType.Field(i).Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				continue
			}
			field := dataType.Field(i)
			fields[name] = field.Type.String()
		}
	}
	return
}

// Kubernetes only allows lower case letters for names.
//
// TODO(ericchiang): refactor ID creation onto the storage.
var encoding = base32.NewEncoding("abcdefghijklmnopqrstuvwxyz234567")

// Valid characters for user codes
const validUserCharacters = "BCDFGHJKLMNPQRSTVWXZ"

// NewDeviceCode returns a 32 char alphanumeric cryptographically secure string
func NewDeviceCode() string {
	return NewSecureID(32)
}

// NewID returns a random string which can be used as an ID for objects.
func NewID() string {
	return NewSecureID(16)
}

func NewSecureID(len int) string {
	buff := make([]byte, len) // random ID.
	if _, err := io.ReadFull(rand.Reader, buff); err != nil {
		panic(err)
	}
	// Avoid the identifier to begin with number and trim padding
	return string(buff[0]%26+'a') + strings.TrimRight(encoding.EncodeToString(buff[1:]), "=")
}

func StringToUint(str string) (int uint) {
	if i, e := strconv.Atoi(str); e == nil {
		return uint(i)
	}
	return 0
}
func StringToInt(str string) int {
	if in, err := strconv.Atoi(str); err != nil {
		return 0
	} else {
		if in < 1 {
			in = 1
		}
		return in
	}
}
func StringToInt64(str string) int64 {
	if in, err := strconv.ParseInt(str, 0, 64); err != nil {
		return 0
	} else {
		if in < 1 {
			in = 1
		}
		return in
	}
}
func StringToFloat64(str string) float64 {
	if in, err := strconv.ParseFloat(str, 64); err != nil {
		return 0
	} else {
		if in < 1 {
			in = 1
		}
		return in
	}
}
func PathExists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// PathIsDir 判断所给路径是否为文件夹
func PathIsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// PathIsFile 判断所给路径是否为文件
func PathIsFile(path string) bool {
	return !PathIsDir(path)
}
func StringInArray(key string, arrays []string) (exist bool) {
	for _, v := range arrays {
		if key == v {
			exist = true
			return
		}
	}
	return
}
func IntInArray(key int, arrays []int) (exist bool) {
	for _, v := range arrays {
		if key == v {
			exist = true
			return
		}
	}
	return
}

func GetRandomString(n int) string {
	randBytes := make([]byte, n/2)
	rand.Read(randBytes)
	return fmt.Sprintf("%x", randBytes)
}
func URL(front, behind string) (url string) {
	if !strings.HasSuffix(front, "/") {
		front += "/"
	}
	if strings.HasPrefix(behind, "/") {
		behind = behind[1:]
	}
	return front + behind
}
