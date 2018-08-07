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
	"github.com/it-chain/engine/icode"
)

type ICodeApi struct {
	ContainerService icode.ContainerService
	GitService       icode.GitService
	EventService     icode.EventService
}

func NewICodeApi(containerService icode.ContainerService, gitService icode.GitService, eventService icode.EventService) ICodeApi {

	return ICodeApi{
		ContainerService: containerService,
		GitService:       gitService,
		EventService:     eventService,
	}
}

func (i ICodeApi) Deploy(id string, baseSaveUrl string, gitUrl string, sshPath string) (icode.Meta, error) {

	// clone meta. in clone function, metaCreatedEvent will publish
	meta, err := i.GitService.Clone(id, baseSaveUrl, gitUrl, sshPath)

	if err != nil {
		return icode.Meta{}, err
	}

	//start ICode with container
	if err = i.ContainerService.StartContainer(meta); err != nil {
		return icode.Meta{}, err
	}

	if err := i.EventService.Publish("meta.created", createMetaCreatedEvent(meta)); err != nil {
		return icode.Meta{}, nil
	}

	return meta, nil
}

func createMetaCreatedEvent(meta icode.Meta) event.MetaCreated {
	return event.MetaCreated{
		ICodeID:        meta.ICodeID,
		Path:           meta.Path,
		Version:        meta.Version,
		CommitHash:     meta.CommitHash,
		GitUrl:         meta.GitUrl,
		RepositoryName: meta.RepositoryName,
	}
}

func (i ICodeApi) UnDeploy(id icode.ID) error {
	// stop iCode container
	err := i.ContainerService.StopContainer(id)

	if err != nil {
		return err
	}

	return i.EventService.Publish("meta.deleted", event.MetaDeleted{ICodeID: id})
}

func (i ICodeApi) ExecuteRequestList(RequestList []icode.Request) []icode.Result {

	resultList := make([]icode.Result, 0)

	for _, request := range RequestList {

		result, err := i.ExecuteRequest(request)

		if err != nil {
			logger.Fatal(nil, fmt.Sprintf("[ICode] request error %s", err.Error()))
		}

		resultList = append(resultList, result)
	}

	return resultList
}

func (i ICodeApi) ExecuteRequest(request icode.Request) (icode.Result, error) {

	return i.ContainerService.ExecuteRequest(request)
}
