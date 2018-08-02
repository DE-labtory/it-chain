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
	"github.com/it-chain/engine/common/rabbitmq/rpc"
)

type BlockApi interface {
	AddBlockToPool(block blockchain.Block) error
	CheckAndSaveBlockFromPool(height blockchain.BlockHeight) error
}

type BlockConfirmCommandHandler struct {
	blockApi BlockApi
}

func NewCommandHandler(blockApi BlockApi) *BlockConfirmCommandHandler {
	return &BlockConfirmCommandHandler{
		blockApi: blockApi,
	}
}

func (h *BlockConfirmCommandHandler) HandleConfirmBlockCommand(command blockchain.ConfirmBlockCommand) (struct{}, rpc.Error) {
	block := command.Block
	if block == nil {
		return struct{}{}, rpc.Error{Message: ErrBlockNil.Error()}
	}

	h.blockApi.AddBlockToPool(block)

	return struct{}{}, rpc.Error{}
}
