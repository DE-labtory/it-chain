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

package tesseract

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/it-chain/engine/ivm"
	"github.com/it-chain/iLogger"
	"github.com/it-chain/tesseract"
	"github.com/it-chain/tesseract/container"
	"github.com/it-chain/tesseract/pb"
	"github.com/rs/xid"
)

var ErrContainerDoesNotExist = errors.New("container does not exist")
var ErrICodeInfoMapNotEmpty = errors.New("ICode info struct in current container is not empty")

type ICodeInfo struct {
	container tesseract.Container
	iCode     ivm.ICode
}

type ContainerService struct {
	sync.RWMutex
	iCodeInfoMap map[tesseract.ContainerID]ICodeInfo
}

func NewContainerService() *ContainerService {
	return &ContainerService{
		iCodeInfoMap: make(map[tesseract.ContainerID]ICodeInfo),
		RWMutex:      sync.RWMutex{},
	}
}

func (cs ContainerService) StartContainer(icode ivm.ICode) error {
	iLogger.Info(nil, fmt.Sprintf("[IVM] Starting container - icodeID: [%s]", icode.ID))
	cs.Lock()
	defer cs.Unlock()

	conf := tesseract.ContainerConfig{
		Name:      icode.RepositoryName,
		Directory: icode.Path,
		Url:       icode.GitUrl,
	}
	container, err := container.Create(conf)

	if err != nil {
		return err
	}

	_, ok := cs.iCodeInfoMap[icode.ID]

	if !ok {
		iCodeInfo := ICodeInfo{
			container: container,
			iCode:     icode,
		}

		cs.iCodeInfoMap[icode.ID] = iCodeInfo
	} else {
		return ErrICodeInfoMapNotEmpty
	}

	return nil
}

func (cs ContainerService) ExecuteRequest(request ivm.Request) (ivm.Result, error) {
	iLogger.Info(nil, fmt.Sprintf("[IVM] Executing icode - icodeID: [%s]", request.ICodeID))

	iCodeInfo, ok := cs.iCodeInfoMap[request.ICodeID]
	if !ok {
		return ivm.Result{}, ErrContainerDoesNotExist
	}

	resultCh := make(chan ivm.Result)
	errCh := make(chan error)

	var callback = func(response *pb.Response, err error) {
		if err != nil {
			errCh <- err
		}

		data := make(map[string]string)
		if len(response.Data) != 0 {
			if err = json.Unmarshal(response.Data, &data); err != nil {
				errCh <- err
				return
			}
		}

		resultCh <- ivm.Result{
			Err:  response.Error,
			Data: data,
		}
	}

	err := iCodeInfo.container.Request(tesseract.Request{
		Uuid:     xid.New().String(),
		Args:     request.Args,
		FuncName: request.Function,
		TypeName: request.Type,
	}, callback)

	if err != nil {
		errCh <- err
	}

	select {
	case err := <-errCh:
		iLogger.Error(nil, fmt.Sprintf("[IVM] fail executing ivm, id:%s", request.ICodeID))
		return ivm.Result{}, err
	case result := <-resultCh:
		return result, nil
	}
}

func (cs ContainerService) StopContainer(id ivm.ID) error {
	iCodeInfo, ok := cs.iCodeInfoMap[id]

	if !ok {
		return ErrContainerDoesNotExist
	}

	err := iCodeInfo.container.Close()

	if err != nil {
		return err
	}

	delete(cs.iCodeInfoMap, id)
	return nil
}

func (cs ContainerService) GetRunningICodeList() []ivm.ICode {
	iCodeList := make([]ivm.ICode, 0)

	for id := range cs.iCodeInfoMap {
		iCodeInfo, ok := cs.iCodeInfoMap[id]
		if ok {
			iCodeList = append(iCodeList, iCodeInfo.iCode)
		}
	}

	return iCodeList
}
