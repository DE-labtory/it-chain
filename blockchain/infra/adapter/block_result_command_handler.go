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

package adapter

import (
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/pkg/errors"
)

var ErrBlockIdNil = errors.New("Error command model ID is nil")
var ErrTxResultsLengthOfZero = errors.New("Error length of tx results is zero")
var ErrTxResultsFail = errors.New("Error not all tx results success")

type BlockResultCommandHandler struct {
	blockQueryService blockchain.BlockQueryService
}

func NewBlockResultCommandHandler(blockQueryService blockchain.BlockQueryService) *BlockResultCommandHandler {
	return &BlockResultCommandHandler{
		blockQueryService: blockQueryService,
	}
}

func (h *BlockResultCommandHandler) HandleBlockResultCommand(command command.ReturnBlockResult) (struct{}, rpc.Error) {
	blockId := command.CommandModel.ID
	txResultList := command.TxResultList

	if blockId == "" {
		return struct{}{}, rpc.Error{Message: ErrBlockIdNil.Error()}
	}

	if len(txResultList) == 0 {
		return struct{}{}, rpc.Error{Message: ErrTxResultsLengthOfZero.Error()}
	}

	if !txResultsSuccess(txResultList) {
		return struct{}{}, rpc.Error{Message: ErrTxResultsFail.Error()}
	}

	err := h.commitBlock(blockId)
	if err != nil {
		return struct{}{}, rpc.Error{Message: err.Error()}
	}

	return struct{}{}, rpc.Error{}
}

func txResultsSuccess(txResultList []command.TxResult) bool {
	allSuccess := true

	for _, result := range txResultList {
		allSuccess = allSuccess && result.Success
	}

	return allSuccess
}

func (h *BlockResultCommandHandler) commitBlock(blockId string) error {
	block, err := h.blockQueryService.GetStagedBlockById(blockId)
	if err != nil {
		return err
	}

	err = blockchain.CommitBlock(&block)
	if err != nil {
		return err
	}

	return nil
}
