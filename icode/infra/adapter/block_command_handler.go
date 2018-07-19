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
	"encoding/json"
	"fmt"
	"sync"

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

func (b *BlockCommandHandler) HandleBlockExecuteCommand(command icode.BlockExecuteCommand) {
	var block icode.Block
	err := json.Unmarshal(command.Block, &block)
	if err != nil {
		fmt.Println("error in handle block excute command. unmashal err")
		return
	}
	b.mutex.Lock()
	results := make([]icode.Result, 0)
	for _, tx := range block.TxList {
		switch tx.TxData.Method {
		case icode.Query:
			results = append(results, *b.icodeApi.Query(tx))
		case icode.Invoke:
			results = append(results, *b.icodeApi.Invoke(tx))
		default:
			fmt.Println(fmt.Sprintf("unknown tx method [%s]", tx.TxData.Method))
		}
	}
	b.commandService.SendBlockExecuteResultCommand(results, command.GetID())
	b.mutex.Unlock()
}
