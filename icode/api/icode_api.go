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
	"github.com/it-chain/engine/icode"
)

type ICodeApi struct {
	ContainerService icode.ContainerService
	GitService       icode.GitService
}

func NewIcodeApi(containerService icode.ContainerService, gitService icode.GitService) *ICodeApi {

	return &ICodeApi{
		ContainerService: containerService,
		GitService:       gitService,
	}
}

func (iApi ICodeApi) Deploy(id string, baseSaveUrl string, gitUrl string, sshPath string) (*icode.Meta, error) {

	// clone meta. in clone function, metaCreatedEvent will publish
	meta, err := iApi.GitService.Clone(id, baseSaveUrl, gitUrl, sshPath)

	if err != nil {
		return nil, err
	}

	//start ICode with container
	if err = iApi.ContainerService.StartContainer(*meta); err != nil {
		return nil, err
	}

	return meta, nil
}

func (iApi ICodeApi) UnDeploy(id icode.ID) error {
	// stop iCode container
	err := iApi.ContainerService.StopContainer(id)
	if err != nil {
		return err
	}

	// publish meta delete event
	err = icode.DeleteMeta(id)
	if err != nil {
		return err
	}

	return nil
}

//todo need asnyc process
func (iApi ICodeApi) Invoke(tx icode.Transaction) *icode.Result {

	result, err := iApi.ContainerService.ExecuteTransaction(tx)

	if err != nil {
		result = &icode.Result{
			TxId:    tx.TxId,
			Data:    nil,
			Success: false,
		}
	}

	return result
}

func (iApi ICodeApi) Query(icodeId icode.ID, functionName string, args []string) *icode.Result {

	tx := icode.Transaction{
		Method:   "query",
		ICodeID:  icodeId,
		Function: functionName,
		Args:     args,
	}

	result, err := iApi.ContainerService.ExecuteTransaction(tx)

	if err != nil {

		result = &icode.Result{
			TxId:    tx.TxId,
			Data:    nil,
			Success: false,
		}
	}
	return result
}
