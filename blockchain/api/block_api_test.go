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

package api_test

import (
	"testing"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/api"
	"github.com/magiconair/properties/assert"
)

func TestBlockApi_AddBlockToPool(t *testing.T) {
	tests := map[string]struct {
		input struct {
			block blockchain.Block
		}
	}{
		"success": {
			input: struct {
				block blockchain.Block
			}{block: &blockchain.DefaultBlock{
				Height: uint64(11),
			}},
		},
	}

	publisherId := "zf"

	blockApi, _ := api.NewBlockApi(publisherId)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		blockApi.AddBlockToPool(test.input.block)
	}
}

func TestBlockApi_CheckAndSaveBlockFromPool(t *testing.T) {
	tests := map[string]struct {
		input struct {
			height blockchain.BlockHeight
		}
		err error
	}{
		"success": {
			input: struct {
				height blockchain.BlockHeight
			}{height: blockchain.BlockHeight(12)},
			err: nil,
		},
	}
	publisherId := "zf"

	// When
	blockApi, _ := api.NewBlockApi(publisherId)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		// When
		err := blockApi.CheckAndSaveBlockFromPool(test.input.height)

		// Then
		assert.Equal(t, test.err, err)
	}
}

func TestBlockApi_SyncIsProgressing(t *testing.T) {
	// when
	publisherId := "zf"

	// when
	blockApi, _ := api.NewBlockApi(publisherId)

	// then
	state := blockApi.SyncIsProgressing()
	assert.Equal(t, blockchain.DONE, state)
}
