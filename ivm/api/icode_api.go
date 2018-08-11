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

package api

import (
	"fmt"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/ivm"
)

type ICodeApi struct {
	ContainerService ivm.ContainerService
	GitService       ivm.GitService
	EventService     ivm.EventService
}

func NewICodeApi(containerService ivm.ContainerService, gitService ivm.GitService, eventService ivm.EventService) ICodeApi {

	return ICodeApi{
		ContainerService: containerService,
		GitService:       gitService,
		EventService:     eventService,
	}
}

func (i ICodeApi) Deploy(id string, baseSaveUrl string, gitUrl string, sshPath string) (ivm.Meta, error) {

	// clone meta. in clone function, metaCreatedEvent will publish
	meta, err := i.GitService.Clone(id, baseSaveUrl, gitUrl, sshPath)

	if err != nil {
		return ivm.Meta{}, err
	}

	//start ICode with container
	if err = i.ContainerService.StartContainer(meta); err != nil {
		return ivm.Meta{}, err
	}

	if err := i.EventService.Publish("meta.created", createMetaCreatedEvent(meta)); err != nil {
		return ivm.Meta{}, nil
	}

	return meta, nil
}

func createMetaCreatedEvent(meta ivm.Meta) event.MetaCreated {
	return event.MetaCreated{
		ICodeID:        meta.ICodeID,
		Path:           meta.Path,
		Version:        meta.Version,
		CommitHash:     meta.CommitHash,
		GitUrl:         meta.GitUrl,
		RepositoryName: meta.RepositoryName,
	}
}

func (i ICodeApi) UnDeploy(id ivm.ID) error {
	// stop iCode container
	err := i.ContainerService.StopContainer(id)

	if err != nil {
		return err
	}

	return i.EventService.Publish("meta.deleted", event.MetaDeleted{ICodeID: id})
}

func (i ICodeApi) ExecuteRequestList(RequestList []ivm.Request) []ivm.Result {

	resultList := make([]ivm.Result, 0)

	for _, request := range RequestList {

		result, err := i.ExecuteRequest(request)

		if err != nil {
			logger.Fatal(nil, fmt.Sprintf("[ICode] request error %s", err.Error()))
		}

		resultList = append(resultList, result)
	}

	return resultList
}

func (i ICodeApi) ExecuteRequest(request ivm.Request) (ivm.Result, error) {

	return i.ContainerService.ExecuteRequest(request)
}
