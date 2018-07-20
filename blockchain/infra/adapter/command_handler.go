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
	"errors"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common/command"
)

var ErrBlockNil = errors.New("Block nil error")

type BlockApi interface {
	AddBlockToPool(block blockchain.Block) error
	CheckAndSaveBlockFromPool(height blockchain.BlockHeight) error
}

type CommandHandler struct {
	blockApi BlockApi
}

func NewCommandHandler(blockApi BlockApi) *CommandHandler {
	return &CommandHandler{
		blockApi: blockApi,
	}
}

// txpool에서 받은 transactions들을 block으로 만들어서 consensus에 보내준다.
func (h *CommandHandler) HandleProposeBlockCommand(command command.ProposeBlock) {
	//rawTxList := command.Transactions
	//
	//txList, err := convertTxList(rawTxList)
	//if err != nil {
	//	// TODO: handle errors
	//	return
	//}
	//
	//block, err := handler.blockApi.CreateBlock(txList)
	//if err != nil {
	//	// TODO: handle errors
	//	return
	//}
	// TODO: service는 api에서 호출되어야한다.
	//dispatcher.SendBlockValidateCommand(block)
}

/// 합의된 block이 넘어오면 block pool에 저장한다.
func (h *CommandHandler) HandleConfirmBlockCommand(command blockchain.ConfirmBlockCommand) error {
	block := command.Block
	if block == nil {
		return ErrBlockNil
	}

	h.blockApi.AddBlockToPool(block)

	return nil
}
