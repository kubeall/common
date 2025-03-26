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

package license

import (
	"encoding/json"
	"fmt"
	"github.com/denisbrodbeck/machineid"
	"github.com/efucloud/common"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"io"
	"net/http"
	"os"
	"path"
	"runtime"
	"time"
)

const (
	k8sPath              = "/var/run/secrets/kubernetes.io/serviceaccount"
	kubernetesServerAddr = "KUBERNETES_PORT_443_TCP_ADDR"
	kubernetesServerPort = "KUBERNETES_SERVICE_PORT"
	dockerEnv            = "/.dockerenv"
)

// GetMachineInformation 根据部署来生成机器信息
func GetMachineInformation(appName string, logger *zap.SugaredLogger) (applicationInfo common.ApplicationInfo) {
	var (
		ca  []byte
		err error
	)
	applicationInfo.Application = appName
	applicationInfo.OS = runtime.GOOS
	applicationInfo.Arch = runtime.GOARCH
	applicationInfo.CpuCores = runtime.GOMAXPROCS(0)
	applicationInfo.Time = time.Now().Local()

	// 判断是否在k8s集群中运行
	ca, err = os.ReadFile(path.Join(k8sPath, "ca.crt"))
	if err == nil {
		applicationInfo.KubernetesInfo = new(common.KubernetesInfo)
		applicationInfo.KubernetesInfo.Version = new(common.K8sVersion)
		applicationInfo.KubernetesInfo.CA = common.MD5VByte(ca)
		tP := path.Join(k8sPath, "namespace")
		if ns, err := os.ReadFile(tP); err == nil {
			applicationInfo.KubernetesInfo.Namespace = string(ns)
		} else {
			applicationInfo.Error = err.Error()
			logger.Errorf("read token from path: %s failed, err: %s", tP, err.Error())
			return
		}
		var (
			k8sTokenPayload *common.K8sTokenPayload
			tokenStr        string
		)
		if token, err := os.ReadFile(path.Join(k8sPath, "token")); err == nil {
			tokenStr = string(token)
			tokenIns, _ := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
				return nil, nil
			})
			data, _ := json.Marshal(tokenIns)
			if json.Unmarshal(data, k8sTokenPayload) == nil {
				if k8sTokenPayload != nil && k8sTokenPayload.Claims != nil && k8sTokenPayload.Claims.KubernetesIo != nil {
					applicationInfo.KubernetesInfo.Namespace = k8sTokenPayload.Claims.KubernetesIo.Namespace
				}
			}
			applicationInfo.MachineID = common.MD5VByte(ca)
		} else {
			logger.Errorf("read token from path: %s failed, err: %s", path.Join(k8sPath, "token"), err.Error())
			applicationInfo.Error = err.Error()
			return
		}
		applicationInfo.KubernetesInfo.Server = os.Getenv(kubernetesServerAddr)
		applicationInfo.KubernetesInfo.Port = os.Getenv(kubernetesServerPort)
		//获取k8s版本信息
		verAddr := fmt.Sprintf("https://%s:%s/version", applicationInfo.KubernetesInfo.Server, applicationInfo.KubernetesInfo.Port)
		logger.Infof("get kubernetes version from: %s", verAddr)
		headers := make(map[string]string)
		headers["Authorization"] = "Bearer " + tokenStr
		if response, err := common.Request(http.MethodGet, verAddr, headers, nil, nil); err == nil {
			body, err := io.ReadAll(response.Body)
			logger.Info(err)
			logger.Infof("get kubernetes version response: %s", string(body))
			if response.StatusCode == http.StatusOK {
				var ver common.K8sVersion
				err = json.Unmarshal(body, &ver)
				if err != nil {
					logger.Error(err)
					applicationInfo.Error = err.Error()
					return
				} else {
					applicationInfo.KubernetesInfo.Version = &ver

				}
			} else {
				logger.Errorf("get kubernetes version response: %s", string(body))
			}
		} else {
			logger.Error(err)
			applicationInfo.Error = err.Error()
			return
		}
	} else {
		logger.Infof("current run system is: %s", runtime.GOOS)
		// 只判断为linux时判断是否docker运行
		if runtime.GOOS == "linux" {
			//只要是linux就认为是容器内部
			if common.PathExists(dockerEnv) {
				applicationInfo.Error = "application not support running in docker"
			}
		} else {
			applicationInfo.PhysicalInfo = new(common.PhysicalInfo)
			applicationInfo.PhysicalInfo.MachineID, _ = machineid.ProtectedID(appName)
		}
	}

	return
}
