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
						PeerID:    "p01",
						ID:        "tx01",
						Status:    0,
						Timestamp: time.Now(),
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
					}},
				},
			},
			err: nil,
		},
	}

	publisher := func(topic string, data interface{}) error {
		command := data.(command.ExecuteBlock)

		assert.Equal(t, topic, "block.execute")
		assert.Equal(t, []byte{0x1}, command.Seal)
		assert.Equal(t, []byte{0x2}, command.PrevSeal)
		assert.Equal(t, blockchain.BlockHeight(1), command.Height)

		return nil
	}

	blockExecuteService := adapter.NewBlockExecuteService(publisher)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		err := blockExecuteService.ExecuteBlock(test.input.block)

		assert.Equal(t, err, test.err)
	}
}
