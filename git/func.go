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

import (
	"bufio"
	"fmt"
	"github.com/efucloud/common"
	"github.com/ghodss/yaml"
	"gopkg.in/ini.v1"
	"io/ioutil"
	"k8s.io/klog/v2"
	"os"
	"path"
	"strings"
	"time"
)

func GetGitRepoInformation() (result *RepositoryInformation) {

	// 从.git/HEAD中读取当前的分支
	curDir, err := os.Getwd()
	if err != nil {
		klog.Errorf("get current work dir failed, err; %s", err.Error())
		return
	}
	gitPath := path.Join(curDir, ".git")
	if !common.PathExists(gitPath) {
		klog.Errorf("%s not exist", gitPath)
		return
	}
	result = &RepositoryInformation{}
	result.Branch, result.Ref, err = getCurrentBranch(path.Join(gitPath, "HEAD"))
	if err != nil {
		return
	}
	result.Hash, err = getCommitHash(path.Join(gitPath, result.Ref))
	if err != nil {
		return
	}
	result.CommitInfo, result.Author, result.Email, result.Commit, result.Time, result.Timestamp, err =
		getCommitLog(path.Join(gitPath, "logs", result.Ref), result.Hash)
	result.Remote, result.Url = getRemoteAndUrl(path.Join(gitPath, "config"), result.Branch)
	return result

}
func getRemoteAndUrl(p, branch string) (remote, url string) {
	cfg, err := ini.Load(p)
	if err != nil {
		klog.Errorf("load file: %s failed, err: %s", p, err.Error())
		return
	}
	sn := fmt.Sprintf(`branch "%s"`, branch)
	section, err := cfg.GetSection(sn)
	if err != nil {
		klog.Errorf("get section: %s failed, err: %s", sn, err.Error())
		return
	}
	key, err := section.GetKey("remote")
	if err != nil {
		klog.Errorf("get section: %s key: remote failed, err: %s", sn, err.Error())
		return
	}
	remote = key.Value()
	remoteSn := fmt.Sprintf(`remote "%s"`, remote)
	section, err = cfg.GetSection(remoteSn)
	if err != nil {
		klog.Errorf("get section: %s failed, err: %s", remoteSn, err.Error())
		return
	}
	key, err = section.GetKey("url")
	if err != nil {
		klog.Errorf("get section: %s key: remote failed, err: %s", sn, err.Error())
		return
	}
	url = key.Value()
	return
}
func getCommitLog(logPath, hash string) (all, author, email, commit, t string, timestamp int64, err error) {
	file, err := os.Open(logPath)
	if err != nil {
		klog.Errorf("read: %s failed, err: %s", logPath, err.Error())
		return
	}
	br := bufio.NewReader(file)
	targetLine := ""
	for {
		l, e := br.ReadBytes('\n')

		if e != nil && len(l) == 0 {
			break
		}
		if strings.Contains(string(l), hash) {
			targetLine = string(l)
			break
		}
	}
	if len(targetLine) > 0 {
		sp := strings.Split(targetLine, hash)
		if len(sp) > 1 {
			all = strings.TrimSpace(strings.Join(sp[1:], ""))
			sps := strings.Split(all, " ")
			if len(sps) >= 4 {
				author = sps[0]
				email = strings.TrimSuffix(strings.TrimPrefix(sps[1], "<"), ">")
				timestamp = common.StringToInt64(sps[2])
				if timestamp > 0 {
					t = time.Unix(timestamp, 0).Format(time.RFC3339)
				}
			}
		}
		commit = strings.TrimSuffix(strings.Trim(strings.Split(targetLine, "commit: ")[1], `"`), "\n")
	}

	return
}

func getCommitHash(p string) (h string, err error) {
	data, err := ioutil.ReadFile(p)
	if err != nil {
		klog.Errorf("read: %s failed, err: %s", p, err.Error())
		return h, err
	}
	return strings.Trim(string(data), "\n"), err
}
func getCurrentBranch(headPath string) (ref, branch string, err error) {
	data, err := ioutil.ReadFile(headPath)
	if err != nil {
		klog.Errorf("read: %s failed, err: %s", headPath, err.Error())
		return ref, headPath, err
	}
	var refD refDef
	err = yaml.Unmarshal(data, &refD)
	if err != nil {
		klog.Errorf("yaml decode: %s failed, err: %s", string(data), err.Error())
		return ref, headPath, err
	}
	b := strings.Split(refD.Ref, "/")
	return b[len(b)-1], refD.Ref, nil
}

type refDef struct {
	Ref string `json:"ref" yaml:"ref"`
}
