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

package adapter_test

import (
	"testing"

	"reflect"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

func TestCommandHandler_HandleConfirmBlockCommand(t *testing.T) {
	tests := map[string]struct {
		input struct {
			command blockchain.ConfirmBlockCommand
		}
		err rpc.Error
	}{
		"success": {
			input: struct {
				command blockchain.ConfirmBlockCommand
			}{
				command: blockchain.ConfirmBlockCommand{
					CommandModel: midgard.CommandModel{ID: "zf"},
					Block: &blockchain.DefaultBlock{
						Height: 99887,
					},
				},
			},
			err: rpc.Error{},
		},
		"block nil error test": {
			input: struct {
				command blockchain.ConfirmBlockCommand
			}{
				command: blockchain.ConfirmBlockCommand{
					CommandModel: midgard.CommandModel{ID: "zf"},
					Block:        nil,
				},
			},
			err: rpc.Error{Message: adapter.ErrBlockNil.Error()},
		},
	}

	blockApi := mock.BlockApi{}
	blockApi.AddBlockToPoolFunc = func(block blockchain.Block) error {
		assert.Equal(t, block.GetHeight(), uint64(99887))
		return nil
	}

	commandHandler := adapter.NewCommandHandler(blockApi)
	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		value, err := commandHandler.HandleConfirmBlockCommand(test.input.command)

		assert.Equal(t, err, test.err)
		assert.True(t, reflect.DeepEqual(value, struct{}{}))
	}

}
