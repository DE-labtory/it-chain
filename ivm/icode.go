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

package ivm

type Version = string
type ID = string
type MetaStatus = int

type ICode struct {
	ID             ID
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        Version
	Language       string
	FolderName     string
}

func NewICode(id string, repositoryName string, folderName string, gitUrl string, path string, commitHash string, version string) ICode {

	return ICode{
		ID:             id,
		CommitHash:     commitHash,
		Path:           path,
		GitUrl:         gitUrl,
		FolderName:     folderName,
		RepositoryName: repositoryName,
		Version:        version,
	}
}
