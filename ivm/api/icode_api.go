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

package api

import (
	"fmt"

	"github.com/DE-labtory/it-chain/common"
	"github.com/DE-labtory/it-chain/common/event"
	"github.com/DE-labtory/it-chain/ivm"
	"github.com/DE-labtory/iLogger"
)

type ICodeApi struct {
	ContainerService ivm.ContainerService
	GitService       ivm.GitService
	EventService     common.EventService
}

func NewICodeApi(containerService ivm.ContainerService, gitService ivm.GitService, eventService common.EventService) ICodeApi {

	return ICodeApi{
		ContainerService: containerService,
		GitService:       gitService,
		EventService:     eventService,
	}
}
func (i ICodeApi) DeployFromRawSsh(baseSaveUrl string, gitUrl string, rawSsh []byte, password string) (ivm.ICode, error) {
	iLogger.Info(nil, fmt.Sprintf("[IVM] Deploying icode - url: [%s]", gitUrl))

	// clone icode. in clone function, metaCreatedEvent will publish
	icode, err := i.GitService.CloneFromRawSsh(baseSaveUrl, gitUrl, rawSsh, password)

	if err != nil {
		return ivm.ICode{}, err
	}

	//on ICode with container
	if err = i.ContainerService.StartContainer(icode); err != nil {
		return ivm.ICode{}, err
	}

	if err := i.EventService.Publish("icode.created", createMetaCreatedEvent(icode)); err != nil {
		return ivm.ICode{}, nil
	}

	iLogger.Info(nil, fmt.Sprintf("[IVM] ICode has deployed - icodeID: [%s]", icode.ID))
	return icode, nil
}

func (i ICodeApi) Deploy(baseSaveUrl string, gitUrl string, sshPath string, password string) (ivm.ICode, error) {
	iLogger.Info(nil, fmt.Sprintf("[IVM] Deploying icode - url: [%s]", gitUrl))

	// clone icode. in clone function, metaCreatedEvent will publish
	icode, err := i.GitService.Clone(baseSaveUrl, gitUrl, sshPath, password)

	if err != nil {
		return ivm.ICode{}, err
	}

	//on ICode with container
	if err = i.ContainerService.StartContainer(icode); err != nil {
		return ivm.ICode{}, err
	}

	if err := i.EventService.Publish("icode.created", createMetaCreatedEvent(icode)); err != nil {
		return ivm.ICode{}, nil
	}

	iLogger.Info(nil, fmt.Sprintf("[IVM] ICode has deployed - icodeID: [%s]", icode.ID))
	return icode, nil
}

func createMetaCreatedEvent(icode ivm.ICode) event.ICodeCreated {
	return event.ICodeCreated{
		ID:             icode.ID,
		Path:           icode.Path,
		Version:        icode.Version,
		CommitHash:     icode.CommitHash,
		GitUrl:         icode.GitUrl,
		RepositoryName: icode.RepositoryName,
	}
}

func (i ICodeApi) UnDeploy(id ivm.ID) error {
	iLogger.Info(nil, fmt.Sprintf("[IVM] Undeploying icode - icodeID: [%s]", id))
	// stop iCode container
	err := i.ContainerService.StopContainer(id)

	if err != nil {
		return err
	}

	iLogger.Info(nil, fmt.Sprintf("[IVM] Icode has undeployed - icodeID: [%s] ", id))

	return i.EventService.Publish("icode.deleted", event.ICodeDeleted{ICodeID: id})
}

func (i ICodeApi) ExecuteRequestList(RequestList []ivm.Request) []ivm.Result {

	resultList := make([]ivm.Result, 0)

	for _, request := range RequestList {

		result, err := i.ExecuteRequest(request)

		if err != nil {
			iLogger.Error(nil, fmt.Sprintf("[IVM] Fail to invoke icode - message: [%s] ", err.Error()))
			result = ivm.Result{Err: err.Error()}
		}

		resultList = append(resultList, result)
	}

	return resultList
}

func (i ICodeApi) ExecuteRequest(request ivm.Request) (ivm.Result, error) {
	return i.ContainerService.ExecuteRequest(request)
}

func (i ICodeApi) GetRunningICodeList() []ivm.ICode {
	return i.ContainerService.GetRunningICodeList()
}
