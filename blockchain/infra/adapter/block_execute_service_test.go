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
	"time"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/magiconair/properties/assert"
)

func TestBlockExecuteService_ExecuteBlock(t *testing.T) {
	tests := map[string]struct {
		input struct {
			block blockchain.Block
		}
		err error
	}{
		"success": {
			input: struct {
				block blockchain.Block
			}{
				block: &blockchain.DefaultBlock{
					Seal:     []byte{0x1},
					PrevSeal: []byte{0x2},
					Height:   blockchain.BlockHeight(1),
					TxList: []*blockchain.DefaultTransaction{{
						ID:        "tx01",
						Status:    0,
						PeerID:    "p01",
						ICodeID:   "icode01",
						Timestamp: time.Now().Round(0),
						TxData: blockchain.TxData{
							Jsonrpc: "jsonRPC01",
							Method:  "invoke",
							Params: blockchain.Params{
								Type:     0,
								Function: "function01",
								Args:     []string{"arg1", "arg2"},
							},
							ID: "txdata01",
						},
						Signature: nil,
					}},
				},
			},
			err: nil,
		},
	}

	server := rpc.NewServer("")
	defer server.Close()

	server.Register("block.execute", func(command command.ExecuteBlock) (command.ExecuteBlock, rpc.Error) {
		return command, rpc.Error{}
	})

	client := rpc.NewClient("")
	defer client.Close()

	blockExecuteService := adapter.NewBlockExecuteService(client)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := blockExecuteService.ExecuteBlock(test.input.block)

		assert.Equal(t, err, test.err)
	}
}
