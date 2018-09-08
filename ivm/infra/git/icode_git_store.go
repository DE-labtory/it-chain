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

	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/ivm"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

var ErrUnsupportedUrl = errors.New("unsupported url [format: github.com/xxx/yyyy], currently only github url is supported")

const (
	github = "github.com"
	gitlab = "gitlab.com"
)

type RepositoryService struct {
}

func NewRepositoryService() *RepositoryService {
	return &RepositoryService{}
}

func (gApi *RepositoryService) Clone(id string, baseSavePath string, repositoryUrl string, sshPath string, password string) (ivm.ICode, error) {
	logger.Info(nil, fmt.Sprintf("[IVM] Cloning Icode - url: [%s]", repositoryUrl))

	giturl, err := toSshUrl(repositoryUrl)
	if err != nil {
		return ivm.ICode{}, err
	}

	name := getNameFromGitUrl(giturl)

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

	sshAuth, err := ssh.NewPublicKeysFromFile("git", sshPath, password)
	if err != nil {
		return ivm.ICode{}, err
	}

	r, err := git.PlainClone(baseSavePath+"/"+name, false, &git.CloneOptions{
		URL:               giturl,
		Auth:              sshAuth,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		return ivm.ICode{}, err
	}

	head, err := r.Head()

	if err != nil {
		return ivm.ICode{}, err
	}

	lastHeadCommit, err := r.CommitObject(head.Hash())
	commitHash := lastHeadCommit.Hash.String()

	if err != nil {
		return ivm.ICode{}, err
	}

	version := "1.0"
	tags, err := r.Tags()
	tags.ForEach(func(tag *plumbing.Reference) error {
		if strings.Compare(tag.Hash().String(), head.Hash().String()) == 0 {
			s := strings.Split(tag.Name().String(), "/")
			version = s[len(s)-1]
		}
		return nil
	})

	logger.Info(nil, fmt.Sprintf("[IVM] ICode has successfully cloned - url: [%s], icodeID: [%s], version[%s]", repositoryUrl, id, version))

	metaData := ivm.NewICode(id, name, repositoryUrl, baseSavePath+"/"+name, commitHash, version)
	return metaData, nil
}

// transfer github.com/it-chain/engine to // git@github.com:it-chain/engine.git
func toSshUrl(repositoryUrl string) (string, error) {
	prefix := "git@"
	postfix := ".git"

	if strings.HasPrefix(repositoryUrl, github) {
		return prefix + github + ":" + after(repositoryUrl, github+"/") + postfix, nil
	}

	return "", ErrUnsupportedUrl
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
