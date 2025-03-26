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
	"time"
)

type K8sTokenClaims struct {
	Aud          []string                    `json:"aud"`
	Exp          int                         `json:"exp"`
	Iat          int                         `json:"iat"`
	Iss          string                      `json:"iss"`
	KubernetesIo *K8sTokenClaimsKubernetesIo `json:"kubernetes.io"`
	Nbf          int                         `json:"nbf"`
	Sub          string                      `json:"sub"`
}
type K8sTokenClaimsKubernetesIo struct {
	Namespace string `json:"namespace"`
	Sub       string `json:"sub"`
}
type K8sTokenPayload struct {
	Claims *K8sTokenClaims `json:"Claims"`
}
type ApplicationPublicInfo struct {
	Application string `json:"application" description:"应用名称"`
	GoVersion   string `json:"goVersion" description:"构建的Go版本"`
	Commit      string `json:"commit" description:"当前构建的Commit"`
	BuildDate   string `json:"buildDate" description:"构建时间"`
	Edition     string `json:"edition" description:"当前版本"`
	BuiltInOrg  string `json:"builtInOrg" description:"内建组织编码"`
	JoinedOrg   bool   `json:"joinedOrg" description:"是否开启加入组织功能"`
}
type ApplicationInfo struct {
	Application    string            `json:"application"`
	GoVersion      string            `json:"goVersion"`
	Commit         string            `json:"commit"`
	BuildDate      string            `json:"buildDate"`
	KubernetesInfo *KubernetesInfo   `json:"kubernetesInfo,omitempty"`
	OS             string            `json:"os,omitempty"`
	Arch           string            `json:"arch,omitempty"`
	CpuCores       int               `json:"cpuCores,omitempty"`
	PhysicalInfo   *PhysicalInfo     `json:"physicalInfo,omitempty"`
	Alert          string            `json:"alert,omitempty"`
	Error          string            `json:"error,omitempty"`
	Time           time.Time         `json:"time,omitempty"`
	Data           string            `json:"data,omitempty"`
	Extend         map[string]string `json:"extend,omitempty"`
	Developer      string            `json:"developer,omitempty"` //
	MachineID      string            `json:"machineId,omitempty"` //
}

type PhysicalInfo struct {
	MachineID  string `json:"machineId"`
	ServerPort string `json:"serverPort"`
}
type KubernetesInfo struct {
	CA        string      `json:"ca,omitempty"`
	Namespace string      `json:"namespace"`
	Server    string      `json:"server"`
	Port      string      `json:"port"`
	Version   *K8sVersion `json:"version"`
}
type K8sVersion struct {
	Major        string    `json:"major"`
	Minor        string    `json:"minor"`
	GitVersion   string    `json:"gitVersion"`
	GitCommit    string    `json:"gitCommit"`
	GitTreeState string    `json:"gitTreeState"`
	BuildDate    time.Time `json:"buildDate"`
	GoVersion    string    `json:"goVersion"`
	Compiler     string    `json:"compiler"`
	Platform     string    `json:"platform"`
}
type Payload struct {
	Data string `json:"data" description:"信息"`
}
