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
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/icode"
	"github.com/it-chain/midgard"
)

type Publisher func(exchange string, topic string, data interface{}) (err error)

type CommandService struct {
	publisher Publisher
}

func NewCommandService(publisher Publisher) *CommandService {
	return &CommandService{
		publisher: publisher,
	}
}

func (c *CommandService) SendBlockExecuteResultCommand(results []icode.Result, blockId string) error {
	return c.publisher("Command", "blockResult", command.BlockResult{
		CommandModel: midgard.CommandModel{
			ID: blockId,
		},
		TxResults: convertTxResults(results),
	})
}

func convertTxResults(icodeResults []icode.Result) []struct {
	TxId    string
	Data    map[string]string
	Success bool
} {
	results := make([]struct {
		TxId    string
		Data    map[string]string
		Success bool
	}, 0)

	for _, result := range icodeResults {
		results = append(results, struct {
			TxId    string
			Data    map[string]string
			Success bool
		}{TxId: result.TxId, Data: result.Data, Success: result.Success})
	}

	return results
}
