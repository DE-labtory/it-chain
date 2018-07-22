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
	"fmt"
	"sync"

	"log"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/icode"
	"github.com/it-chain/engine/icode/api"
)

type BlockCommandHandler struct {
	icodeApi       api.ICodeApi
	commandService icode.CommandService
	mutex          *sync.Mutex
}

func NewBlockCommandHandler(icodeApi api.ICodeApi, service icode.CommandService) *BlockCommandHandler {
	return &BlockCommandHandler{
		icodeApi:       icodeApi,
		commandService: service,
		mutex:          &sync.Mutex{},
	}
}

func (b *BlockCommandHandler) HandleBlockExecuteCommand(blockExecuteCommand command.ExecuteBlock) {

	b.mutex.Lock()
	defer b.mutex.Unlock()

	results := make([]icode.Result, 0)
	for _, tx := range blockExecuteCommand.TxList {
		switch tx.Method {
		case icode.Query:
			log.Print(fmt.Sprintf("dont excute query using blockExecuteCommand , txID : [%s]", tx.ID))
		case icode.Invoke:
			results = append(results, *b.icodeApi.Invoke(icode.Transaction{
				TxId:     tx.ID,
				ICodeID:  tx.ICodeID,
				Function: tx.Function,
				Method:   tx.Method,
				Jsonrpc:  tx.Jsonrpc,
				Args:     tx.Args,
			}))
		default:
			fmt.Println(fmt.Sprintf("unknown tx method [%s]", tx.Method))
		}
	}

	b.commandService.SendBlockExecuteResultCommand(results, blockExecuteCommand.GetID())
}
