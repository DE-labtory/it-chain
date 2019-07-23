/*
 * Copyright 2018 DE-labtory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package git

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/DE-labtory/it-chain/ivm"
	"github.com/DE-labtory/iLogger"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var ErrUnsupportedUrl = errors.New("unsupported url [format: github.com/xxx/yyyy], currently github and gitlab url is supported")

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

const (
	github         = "github.com"
	gitlab         = "gitlab.com"
	defaultVersion = "1.0"
)

type RepositoryService struct {
}

func NewRepositoryService() *RepositoryService {
	return &RepositoryService{}
}

func (gApi *RepositoryService) CloneFromRawSsh(baseSavePath string, repositoryUrl string, rawSsh []byte, password string) (ivm.ICode, error) {
	var randomString = RandStringBytesMaskImprSrc(40)

	iLogger.Infof(nil, "[IVM] Cloning Icode - url: [%s]", repositoryUrl)

	var gitUrl string
	var girUrlErr error

	if len(rawSsh) == 0 {
		gitUrl, girUrlErr = toGitUrl(repositoryUrl, "https")
	} else {
		gitUrl, girUrlErr = toGitUrl(repositoryUrl, "git")
	}

	if girUrlErr != nil {
		return ivm.ICode{}, girUrlErr
	}

	name := getNameFromGitUrl(gitUrl)
	if name == "" {
		return ivm.ICode{}, errors.New(fmt.Sprintf("Invalid url name [%s]", repositoryUrl))
	}

	cloneOptions := &git.CloneOptions{
		URL:               gitUrl,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	if len(rawSsh) != 0 {
		sshAuth, err := ssh.NewPublicKeys("git", rawSsh, password)
		if err != nil {
			return ivm.ICode{}, err
		}
		cloneOptions.Auth = sshAuth

	}

	r, err := git.PlainClone(baseSavePath+"/"+name+"_"+randomString, false, cloneOptions)
	if err != nil {
		return ivm.ICode{}, err
	}

	head, err := r.Head()

	if err != nil {
		return ivm.ICode{}, err
	}

	lastHeadCommit, err := r.CommitObject(head.Hash())
	if err != nil {
		return ivm.ICode{}, err
	}

	commitHash := lastHeadCommit.Hash.String()

	version := defaultVersion
	tags, err := r.Tags()
	if err != nil {
		return ivm.ICode{}, err
	}

	tags.ForEach(func(tag *plumbing.Reference) error {
		if strings.Compare(tag.Hash().String(), head.Hash().String()) == 0 {
			s := strings.Split(tag.Name().String(), "/")
			version = s[len(s)-1]
		}
		return nil
	})

	id := name + "_" + commitHash

	//check file already exist

	iLogger.Info(nil, fmt.Sprintf("[IVM] ICode has successfully cloned - url: [%s], icodeID: [%s], version[%s]", repositoryUrl, id, version))

	metaData := ivm.NewICode(id, name, name+"_"+randomString, repositoryUrl, baseSavePath+"/"+name+"_"+randomString, commitHash, version)
	return metaData, nil
}

func (gApi *RepositoryService) Clone(baseSavePath string, repositoryUrl string, sshPath string, password string) (ivm.ICode, error) {

	if sshPath == "" {
		return gApi.CloneFromRawSsh(baseSavePath, repositoryUrl, make([]byte, 0), password)
	}

	bytes, err := ioutil.ReadFile(sshPath)
	if err != nil {
		return ivm.ICode{}, err
	}

	return gApi.CloneFromRawSsh(baseSavePath, repositoryUrl, bytes, password)
}

// transfer github.com/DE-labtory/it-chain to
// is ssh (git@github.com:DE-labtory/it-chain.git) or is not ssh (https://github.com/DE-labtory/it-chain.git)
func toGitUrl(repositoryUrl string, protocolType string) (string, error) {
	const postfix = ".git"

	var gitUrlsort string
	if strings.HasPrefix(repositoryUrl, github) {
		gitUrlsort = github
	} else if strings.HasPrefix(repositoryUrl, gitlab) {
		gitUrlsort = gitlab
	} else {
		return "", ErrUnsupportedUrl
	}

	switch protocolType {
	case "git":
		return "git@" + gitUrlsort + ":" + after(repositoryUrl, gitUrlsort+"/") + postfix, nil
	case "https":
		return "https://" + gitUrlsort + "/" + after(repositoryUrl, gitUrlsort+"/") + postfix, nil
	default:
		return "", errors.New("unsupported type")
	}
}

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}

	return value[adjustedPos:len(value)]
}

func getNameFromGitUrl(gitUrl string) string {
	parsed := strings.Split(gitUrl, "/")

	if len(parsed) == 0 {
		return ""
	}

	name := strings.Split(parsed[len(parsed)-1], ".")

	if len(name) == 0 {
		return ""
	}

	return name[0]
}
