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
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/stretchr/testify/assert"
	"time"
	"os"
	"io/ioutil"
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
	blockRepo := mock.BlockRepository{}
	eventService := mock.EventService{}

	blockApi, _ := api.NewBlockApi(publisherId, blockRepo, eventService)

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
	blockRepo := mock.BlockRepository{}
	eventService := mock.EventService{}


	// When
	blockApi, _ := api.NewBlockApi(publisherId, blockRepo, eventService)

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
	blockRepo := mock.BlockRepository{}
	eventService := mock.EventService{}

	// when
	blockApi, _ := api.NewBlockApi(publisherId, blockRepo, eventService)

	// then
	state := blockApi.SyncIsProgressing()
	assert.Equal(t, blockchain.DONE, state)
}

// TODO: Write real situation test code, after finishing implementing api_gatey block_query_api.go
func TestBlockApi_CommitProposedBlock(t *testing.T) {

	publisherID := "junksound"

	txList := []*blockchain.DefaultTransaction{
		{
			ID:        "tx01",
			ICodeID:   "ICodeID",
			PeerID:    "junksound",
			Timestamp: time.Now().Round(0),
			Jsonrpc:   "",
			Function:  "",
			Args:      make([]string, 0),
			Signature: []byte("Signature"),
		},
	}

	lastBlock := blockchain.DefaultBlock{
		Seal:      []byte("seal"),
		PrevSeal:  []byte("prevSeal"),
		Height:    uint64(11),
		TxList:    txList,
		TxSeal:    nil,
		Timestamp: time.Time{},
		Creator:   nil,
		State:     "",
	}

	blockRepo := mock.BlockRepository{}
	blockRepo.FindLastFunc = func() (blockchain.DefaultBlock, error) {
		return lastBlock, nil
	}

	blockRepo.SaveFunc = func(block blockchain.DefaultBlock) error {

		assert.Equal(t,[]byte("seal"),block.PrevSeal)

		return nil
	}

	eventService := mock.EventService{}
	eventService.PublishFunc = func(topic string, event interface{}) error {
		return nil
	}


	bApi, err := api.NewBlockApi(publisherID,blockRepo,eventService)
	assert.NoError(t, err)

	// when
	err = bApi.CommitProposedBlock(txList)
	assert.NoError(t, err)
}

func TestBlockApi_CommitGenesisBlock(t *testing.T) {
	GenesisFilePath := "./Genesis.conf"
	defer os.Remove(GenesisFilePath)

	GenesisBlockConfigJson := []byte(`{
									"Orgainaization":"Default",
									"NetworkId":"Default",
								  	"Height":0,
								  	"TimeStamp":"Jan 1, 2018 at 0:00am (KST)",
								  	"Creator":"junksound"
								}`)

	err := ioutil.WriteFile(GenesisFilePath, GenesisBlockConfigJson, 0644)
	assert.NoError(t, err)

	blockRepo := mock.BlockRepository{}
	blockRepo.SaveFunc = func(block blockchain.DefaultBlock) error {
		return nil
	}

	eventService := mock.EventService{}
	eventService.PublishFunc = func(topic string, event interface{}) error {
		return nil
	}

	publisherID := "junksound"
	bApi, err := api.NewBlockApi(publisherID,blockRepo,eventService)
	assert.NoError(t, err)
	// when
	err = bApi.CommitGenesisBlock(GenesisFilePath)
	assert.NoError(t, err)
}