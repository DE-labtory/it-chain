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

func (i ICodeApi) Deploy(id string, baseSaveUrl string, gitUrl string, sshPath string) (*icode.Meta, error) {

	// clone meta. in clone function, metaCreatedEvent will publish
	meta, err := i.GitService.Clone(id, baseSaveUrl, gitUrl, sshPath)

	if err != nil {
		return nil, err
	}

	//start ICode with container
	if err = i.ContainerService.StartContainer(*meta); err != nil {
		return nil, err
	}

	return meta, nil
}

func (i ICodeApi) UnDeploy(id icode.ID) error {
	// stop iCode container
	err := i.ContainerService.StopContainer(id)
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

func (i ICodeApi) ExecuteTransactionList(transactionList []icode.Transaction) []icode.Result {

	resultList := make([]icode.Result, 0)

	for _, transaction := range transactionList {
		result := i.ExecuteTransaction(transaction)
		resultList = append(resultList, result)
	}

	return resultList
}

func (i ICodeApi) ExecuteTransaction(tx icode.Transaction) icode.Result {

	result, err := i.ContainerService.ExecuteTransaction(tx)

	if err != nil {
		result = &icode.Result{
			TxId:    tx.TxId,
			Data:    nil,
			Success: false,
		}
	}

	return *result
}
