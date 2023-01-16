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

package git

type RepositoryInformation struct {
	Remote     string `json:"remote" yaml:"remote"`
	Url        string `json:"url" yaml:"url"`
	Branch     string `json:"branch" yaml:"branch"`
	Ref        string `json:"ref" yaml:"ref"`
	Hash       string `json:"hash" yaml:"hash"`
	Author     string `json:"author" yaml:"author"`
	Email      string `json:"email" yaml:"email"`
	Timestamp  int64  `json:"timestamp" yaml:"timestamp"`
	Time       string `json:"time" yaml:"time"`
	Commit     string `json:"commit" yaml:"commit"`
	CommitInfo string `json:"commitInfo" yaml:"commitInfo"`
}
