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

package eauth

import "github.com/golang-jwt/jwt/v5"

type UserInfo struct {
	Subject          string                 `json:"sub"`
	Profile          string                 `json:"profile"`
	Email            string                 `json:"email"`
	EmailVerified    bool                   `json:"email_verified"`
	Org              string                 `json:"org,omitempty"`
	Customs          map[string]interface{} `json:"customs"` // 组织自定义属性
	Providers        []string               `json:"providers"`
	Groups           []string               `json:"groups"`
	RegistrationFrom string                 `json:"registrationFrom"` // 注册渠道
	AuthProvider     string                 `json:"authProvider"`     // 认证提供商
	Username         string                 `json:"username"`         // 用户名 组织内唯一必须由DNS-1123标签格式的单元组成
	Nickname         string                 `json:"nickname"`         // 昵称，如中文名
	Role             string                 `json:"role,omitempty"`   //组织角色
	Category         string                 `json:"category"`
	Phone            string                 `json:"phone"`
	ID               uint                   `json:"id"`
	Enable           bool                   `json:"enable"`
	Workspaces       []string               `json:"workspaces" yaml:"workspaces"` // 工作空间
	WorkspacesRoles  map[string][]string    `json:"workspacesRoles"`              // 工作空间角色
}

type ApplicationSyncAccountInfo struct {
	EAuthID          uint                   `json:"eAuthId" yaml:"eAuthId"`
	Organization     string                 `json:"organization" validate:"required"` // 组织编码
	Username         string                 `json:"username" validate:"dns1123"`      // 用户名 组织内唯一必须由DNS-1123标签格式的单元组成
	Nickname         string                 `json:"nickname"`                         // 昵称，如中文名
	AdminApps        []string               `json:"adminApps"`                        // 应用管理员
	Enable           uint                   `json:"enable" validate:"oneof=0 1"`      // 是否有效，组织管理员不能设置为无效
	OrgCustoms       map[string]interface{} `json:"orgCustoms"`                       // 组织自定义属性
	Hash             string                 `json:"hash"`                             // 组织:用户名的Hash
	RegistrationFrom string                 `json:"registrationFrom"`                 // 注册渠道
	Language         string                 `json:"language" validate:"oneof=en zh"`  // 语言
	Email            string                 `json:"email" yaml:"email"`
	Phone            string                 `json:"phone" yaml:"phone"`
	Groups           []string               `json:"groups" yaml:"groups"`
	Workspaces       []string               `json:"workspaces" yaml:"workspaces"` // 工作空间
	WorkspacesRoles  map[string][]string    `json:"workspacesRoles"`              // 工作空间角色
}
type AccountClaims struct {
	EAuthID         uint                `json:"eAuthId"`
	Org             string              `json:"org,omitempty"`
	AuthProvider    string              `json:"authProvider"`
	Username        string              `json:"username"` // 用户名 组织内唯一必须由DNS-1123标签格式的单元组成
	Nickname        string              `json:"nickname"` // 昵称，如中文名
	Role            string              `json:"role"`     // 组织角色
	Nonce           string              `json:"nonce"`
	Email           string              `json:"email"`
	Phone           string              `json:"phone"`
	Groups          []string            `json:"groups"`
	Workspaces      []string            `json:"workspaces"`      // 工作空间
	WorkspacesRoles map[string][]string `json:"workspacesRoles"` // 工作空间角色
	AppCode         string              `json:"appCode"`
	AppClientID     string              `json:"appClientId"`
	AppOwner        bool                `json:"appOwner"`
	Category        string              `json:"category"`
	jwt.RegisteredClaims
}

// LocalLoginParam 本地登录请求
type LocalLoginParam struct {
	Method      string `json:"method" validate:"oneof=password phoneCode emailCode"` // 登录类型，用户名密码/手机验证码/邮箱验证码/
	Username    string `json:"username"`
	Password    string `json:"password"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	ValidCode   string `json:"validCode"`
	Code        string `json:"code"`
	State       string `json:"state"`
	RedirectUri string `json:"redirectUri" validate:"required"`
	Bind        string `json:"bind"`
}
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	IDToken      string `json:"id_token"`
}

type RefreshToken struct {
	Org       string `json:"org"`
	App       string `json:"app"`
	ExpiresIn int64  `json:"expiresIn"`
	AccountID uint   `json:"accountId"`
	Issuer    string `json:"issuer"`
	Provider  string `json:"provider"`
}

type AccountSync struct {
	CronJob string `json:"cronJob" yaml:"cronJob"` //
	Address string `json:"address" yaml:"address"` //
}
type WorkspaceSync struct {
	CronJbo    string   `json:"cronJbo" yaml:"cronJbo"`       //
	Address    string   `json:"address" yaml:"address"`       //
	Workspaces []string `json:"workspaces" yaml:"workspaces"` //
}
