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

	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/ivm"
	"github.com/it-chain/tesseract"
	"github.com/it-chain/tesseract/container"
	"github.com/it-chain/tesseract/pb"
	"github.com/rs/xid"
)

var ErrContainerDoesNotExist = errors.New("container does not exist")

type ContainerService struct {
	sync.RWMutex
	containerMap map[tesseract.ContainerID]tesseract.Container
}

func NewContainerService() *ContainerService {

	return &ContainerService{
		containerMap: make(map[tesseract.ContainerID]tesseract.Container),
		RWMutex:      sync.RWMutex{},
	}
}

func (cs ContainerService) StartContainer(icode ivm.ICode) error {
	logger.Info(nil, fmt.Sprintf("[ICode] deploying ivm, id:%s", icode.ID))
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

	cs.containerMap[icode.ID] = container

	return nil
}

func (cs ContainerService) ExecuteRequest(request ivm.Request) (ivm.Result, error) {
	logger.Info(nil, fmt.Sprintf("[ICode] executing ivm, id:%s", request.ICodeID))

	container, ok := cs.containerMap[request.ICodeID]

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

	err := container.Request(tesseract.Request{
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
		logger.Error(nil, fmt.Sprintf("[ICode] fail executing ivm, id:%s", request.ICodeID))
		return ivm.Result{}, err
	case result := <-resultCh:
		return result, nil
	}
}

func (cs ContainerService) StopContainer(id ivm.ID) error {

	container, ok := cs.containerMap[id]

	if !ok {
		return ErrContainerDoesNotExist
	}

	err := container.Close()

	if err != nil {
		return err
	}

	delete(cs.containerMap, id)
	return nil
}

func (cs ContainerService) GetRunningICodeIDList() []ivm.ID {

	icodeIDList := make([]ivm.ID, 0)

	for id, _ := range cs.containerMap {
		icodeIDList = append(icodeIDList, id)
	}

	return icodeIDList
}
