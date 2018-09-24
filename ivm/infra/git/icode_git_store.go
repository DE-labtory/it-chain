/*
 * Copyright 2018 It-chain
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
	"os"
	"strings"

	"github.com/it-chain/engine/ivm"
	"github.com/it-chain/iLogger"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var ErrUnsupportedUrl = errors.New("unsupported url [format: github.com/xxx/yyyy], currently github and gitlab url is supported")

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

func (gApi *RepositoryService) Clone(id string, baseSavePath string, repositoryUrl string, sshPath string, password string) (ivm.ICode, error) {
	iLogger.Info(nil, fmt.Sprintf("[IVM] Cloning Icode - url: [%s]", repositoryUrl))

	gitUrl, err := toGitUrl(repositoryUrl, sshPath)
	if err != nil {
		return ivm.ICode{}, err
	}

	name := getNameFromGitUrl(gitUrl)

	if name == "" {
		return ivm.ICode{}, errors.New(fmt.Sprintf("Invalid url name [%s]", repositoryUrl))
	}

	//check file already exist
	if _, err := os.Stat(baseSavePath + "/" + name); err == nil {
		// if iCode already exist, remove that
		err = os.RemoveAll(baseSavePath + "/" + name)
		if err != nil {
			return ivm.ICode{}, err
		}
	}

	cloneOptions := &git.CloneOptions{
		URL:               gitUrl,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	}

	if sshPath != "" {
		sshAuth, err := ssh.NewPublicKeysFromFile("git", sshPath, password)
		if err != nil {
			return ivm.ICode{}, err
		}

		cloneOptions.Auth = sshAuth
	}

	r, err := git.PlainClone(baseSavePath+"/"+name, false, cloneOptions)

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

	iLogger.Info(nil, fmt.Sprintf("[IVM] ICode has successfully cloned - url: [%s], icodeID: [%s], version[%s]", repositoryUrl, id, version))

	metaData := ivm.NewICode(id, name, repositoryUrl, baseSavePath+"/"+name, commitHash, version)
	return metaData, nil
}

// transfer github.com/it-chain/engine to
// is ssh (git@github.com:it-chain/engine.git) or is not ssh (https://github.com/it-chain/engine.git)
func toGitUrl(repositoryUrl string, sshAuth string) (string, error) {
	const postfix = ".git"

	var gitUrlsort string
	if strings.HasPrefix(repositoryUrl, github) {
		gitUrlsort = github
	} else if strings.HasPrefix(repositoryUrl, gitlab) {
		gitUrlsort = gitlab
	} else {
		return "", ErrUnsupportedUrl
	}

	var prefix string
	if sshAuth != "" {
		prefix = "git@" + gitUrlsort + ":"
	} else {
		prefix = "https://" + gitUrlsort + "/"
	}

	return prefix + after(repositoryUrl, gitUrlsort+"/") + postfix, nil
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
