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

package service

import (
	"errors"
	"fmt"

	"encoding/json"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/engine/icode"
	"github.com/it-chain/midgard"
	"github.com/it-chain/tesseract"
	"github.com/it-chain/tesseract/cellcode/cell"
)

type TesseractContainerService struct {
	tesseract      *tesseract.Tesseract
	containerIdMap map[icode.ID]string // key : iCodeId, value : containerId
}

func NewTesseractContainerService(config tesseract.Config) *TesseractContainerService {
	tesseractObj := &TesseractContainerService{
		tesseract:      tesseract.New(config),
		containerIdMap: make(map[icode.ID]string, 0),
	}
	return tesseractObj
}

func (cs TesseractContainerService) StartContainer(meta icode.Meta) error {

	tesseractIcodeInfo := tesseract.ICodeInfo{
		Name:      meta.RepositoryName,
		Directory: meta.Path,
	}

	containerId, err := cs.tesseract.SetupContainer(tesseractIcodeInfo)

	if err != nil {
		icode.ChangeMetaStatus(meta.GetID(), icode.DEPLOY_FAIL)
		return err
	}

	cs.containerIdMap[meta.ICodeID] = containerId
	icode.ChangeMetaStatus(meta.GetID(), icode.DEPLOYED)

	return nil
}

func (cs TesseractContainerService) ExecuteTransaction(tx icode.Transaction) (*icode.Result, error) {

	containerId, found := cs.containerIdMap[tx.ICodeID]

	if !found {
		return nil, errors.New(fmt.Sprintf("no container for iCode : %s", tx.ICodeID))
	}

	tesseractTxInfo := cell.TxInfo{
		Method: tx.Method,
		ID:     tx.ICodeID,
		Params: cell.Params{
			Function: tx.Function,
			Args:     tx.Args,
		},
	}

	res, err := cs.tesseract.QueryOrInvoke(containerId, tesseractTxInfo)

	if err != nil {
		return nil, err
	}
	var data = make(map[string]string)
	var isSuccess bool

	switch res.Result {
	case "Success":
		isSuccess = true
		err = json.Unmarshal(res.Data, &data)
		if err != nil {
			return nil, err
		}

	case "Error":
		isSuccess = false
		data = nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown pb response result %s", res.Result))
	}

	result := &icode.Result{
		Data:    data,
		TxId:    tx.TxId,
		Success: isSuccess,
	}
	return result, nil
}

func (cs TesseractContainerService) StopContainer(id icode.ID) error {
	containerId := cs.containerIdMap[id]
	if containerId == "" {
		return errors.New(fmt.Sprintf("no container with icode id %s:", id))
	}
	err := cs.tesseract.StopContainerById(containerId)
	if err != nil {
		return err
	}
	delete(cs.containerIdMap, id)
	deletedEvent := event.MetaDeleted{
		EventModel: midgard.EventModel{
			ID:   id,
			Type: "meta.deleted",
		},
	}
	return eventstore.Save(id, deletedEvent)
}
