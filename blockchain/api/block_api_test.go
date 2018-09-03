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

	"encoding/hex"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/api"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/stretchr/testify/assert"
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
	blockPool := blockchain.NewBlockPool()

	blockApi, _ := api.NewBlockApi(publisherId, blockRepo, eventService, blockPool)

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
	blockPool := blockchain.NewBlockPool()

	// When
	blockApi, _ := api.NewBlockApi(publisherId, blockRepo, eventService, blockPool)

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
	blockPool := blockchain.NewBlockPool()

	// when
	blockApi, _ := api.NewBlockApi(publisherId, blockRepo, eventService, blockPool)

	// then
	state := blockApi.SyncIsProgressing()
	assert.Equal(t, blockchain.DONE, state)
}

// TODO: Write real situation test code, after finishing implementing api_gatey block_query_api.go
func TestBlockApi_CommitProposedBlock(t *testing.T) {

	lastBlock := blockchain.DefaultBlock{
		Seal:     []byte("seal"),
		PrevSeal: []byte("prevSeal"),
		Height:   uint64(11),
		TxList: []*blockchain.DefaultTransaction{
			{
				ID:        "lastBlockID",
				ICodeID:   "ICodeID",
				PeerID:    "junksound",
				Timestamp: time.Now().Round(0),
				Jsonrpc:   "",
				Function:  "",
				Args:      make([]string, 0),
				Signature: []byte("Signature"),
			},
		},
		TxSeal:    nil,
		Timestamp: time.Time{},
		Creator:   nil,
		State:     "",
	}

	block := mock.GetNewBlock(lastBlock.GetSeal(), lastBlock.GetHeight()+1)

	//set subscriber
	var wg sync.WaitGroup
	wg.Add(1)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	handler := &mock.CommitEventHandler{}
	handler.HandleFunc = func(event event.BlockCommitted) {
		assert.Equal(t, "tx01", event.TxList[0].ID)
		assert.Equal(t, blockchain.Committed, event.State)
		wg.Done()
	}

	subscriber.SubscribeTopic("block.*", handler)

	publisherID := "junksound"

	blockRepo := mock.BlockRepository{}
	blockRepo.FindLastFunc = func() (blockchain.DefaultBlock, error) {
		return lastBlock, nil
	}
	blockRepo.SaveFunc = func(block blockchain.DefaultBlock) error {
		assert.Equal(t, "tx01", block.GetTxList()[0].GetID())
		return nil
	}

	eventService := common.NewEventService("", "Event")
	blockPool := blockchain.NewBlockPool()

	bApi, err := api.NewBlockApi(publisherID, blockRepo, eventService, blockPool)
	assert.NoError(t, err)
	// when
	err = bApi.CommitBlock(*block)

	// then
	assert.NoError(t, err)
	wg.Wait()
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

	//set subscriber
	var wg sync.WaitGroup
	wg.Add(1)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	handler := &mock.CommitEventHandler{}

	handler.HandleFunc = func(event event.BlockCommitted) {
		assert.Equal(t, hex.EncodeToString([]byte("junksound")), hex.EncodeToString(event.Creator))
		assert.Equal(t, blockchain.Committed, event.State)
		wg.Done()
	}

	subscriber.SubscribeTopic("block.*", handler)

	publisherID := "junksound"

	blockRepo := mock.BlockRepository{}

	blockRepo.SaveFunc = func(block blockchain.DefaultBlock) error {
		assert.Equal(t, hex.EncodeToString([]byte("junksound")), hex.EncodeToString(block.GetCreator()))
		return nil
	}

	eventService := common.NewEventService("", "Event")
	blockPool := blockchain.NewBlockPool()

	bApi, err := api.NewBlockApi(publisherID, blockRepo, eventService, blockPool)
	assert.NoError(t, err)

	// when
	err = bApi.CommitGenesisBlock(GenesisFilePath)
	// then
	assert.NoError(t, err)
	wg.Wait()
}

func TestBlockApi_CreateProposedBlock(t *testing.T) {
	// given
	publisherID := "zf"

	lastBlock := mock.GetNewBlock([]byte("prevSeal"), 1)

	blockRepo := mock.BlockRepository{}
	blockRepo.FindLastFunc = func() (blockchain.DefaultBlock, error) {
		return *lastBlock, nil
	}

	eventService := mock.EventService{}
	blockPool := blockchain.NewBlockPool()

	blockApi, err := api.NewBlockApi(publisherID, blockRepo, eventService, blockPool)
	assert.NoError(t, err)

	txList := mock.GetTxList(time.Now())

	// when
	block, err := blockApi.CreateProposedBlock(txList)

	// then
	assert.NoError(t, err)
	assert.Equal(t, lastBlock.GetSeal(), block.GetPrevSeal())
	assert.Equal(t, uint64(2), block.GetHeight())
}

func TestBlockApi_StageBlock(t *testing.T) {
	// when
	block := &blockchain.DefaultBlock{
		Height: 30,
		State:  blockchain.Committed,
	}
	publisherId := "1"
	blockRepo := mock.BlockRepository{}
	eventService := mock.EventService{}
	blockPool := blockchain.NewBlockPool()

	// when
	blockApi, _ := api.NewBlockApi(publisherId, blockRepo, eventService, blockPool)
	blockApi.StageBlock(*block)

	// when
	toBe := &blockchain.DefaultBlock{
		Height: 30,
		State:  blockchain.Staged,
	}

	// then
	assert.Equal(t, blockApi.BlockPool.Get(30), toBe)
}
