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

	"github.com/it-chain/engine/icode"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

type RepositoryService struct {
}

func NewRepositoryService() *RepositoryService {
	return &RepositoryService{}
}

func (gApi *RepositoryService) Clone(id string, baseSavePath string, repositoryUrl string, sshPath string) (icode.Meta, error) {
	name := getNameFromGitUrl(repositoryUrl)

	if name == "" {
		return icode.Meta{}, errors.New(fmt.Sprintf("Invalid url name [%s]", repositoryUrl))
	}

	//check file already exist
	if _, err := os.Stat(baseSavePath + "/" + name); err == nil {
		// if iCode already exist, remove that
		err = os.RemoveAll(baseSavePath + "/" + name)
		if err != nil {
			return icode.Meta{}, err
		}
	}

	sshAuth, err := ssh.NewPublicKeysFromFile("git", sshPath, "")
	if err != nil {
		return icode.Meta{}, err
	}

	r, err := git.PlainClone(baseSavePath+"/"+name, false, &git.CloneOptions{
		URL:               repositoryUrl,
		Auth:              sshAuth,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})

	if err != nil {
		return icode.Meta{}, err
	}

	head, err := r.Head()

	if err != nil {
		return icode.Meta{}, err
	}

	lastHeadCommit, err := r.CommitObject(head.Hash())
	commitHash := lastHeadCommit.Hash.String()

	if err != nil {
		return icode.Meta{}, err
	}

	metaData := icode.NewMeta(id, name, repositoryUrl, baseSavePath+"/"+name, commitHash)
	return metaData, nil
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
