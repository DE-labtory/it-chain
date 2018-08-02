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
)

type BlockResultCommandHandler struct {
	blockQueryService blockchain.BlockQueryService
}

func NewBlockResultCommandHandler(blockQueryService blockchain.BlockQueryService) *BlockResultCommandHandler {
	return &BlockResultCommandHandler{
		blockQueryService: blockQueryService,
	}
}

func (h *BlockResultCommandHandler) HandleBlockResultCommand(command command.ReturnBlockResult) error {
	blockId := command.CommandModel.ID
	txResultList := command.TxResultList

	if blockId == "" {
		return ErrBlockIdNil
	}

	if len(txResultList) == 0 {
		return ErrTxResultsLengthOfZero
	}

	if !txResultsSuccess(txResultList) {
		return ErrTxResultsFail
	}

	err := h.commitBlock(blockId)
	if err != nil {
		return err
	}

	return nil
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
