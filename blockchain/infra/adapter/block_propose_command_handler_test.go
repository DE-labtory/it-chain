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

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"

	"os"
	"sync"

	"time"

	"errors"

	"github.com/it-chain/engine/blockchain/api"
	"github.com/it-chain/engine/blockchain/infra/mem"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestBlockProposeCommandHandler_HandleProposeBlockCommand_Solo(t *testing.T) {

	//set subscriber
	var wg sync.WaitGroup
	wg.Add(2)

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	handler := &mock.CommitEventHandler{}
	handler.HandleFunc = func(event event.BlockCommitted) {
		assert.Equal(t, "tx01", event.TxList[0].ID)
		assert.Equal(t, blockchain.Committed, event.State)
		wg.Done()
	}

	subscriber.SubscribeTopic("block.*", handler)

	consensusService := mock.ConsensusService{}
	consensusService.ConsensusBlockFunc = func(block blockchain.DefaultBlock) error {
		return nil
	}

	//set bApi
	publisherID := "junksound"
	dbPath := "./.db"

	br, err := mem.NewBlockRepository(dbPath)
	assert.Equal(t, nil, err)
	defer func() {
		br.Close()
		os.RemoveAll(dbPath)
	}()

	prevBlock := mock.GetNewBlock([]byte("genesis"), 0)

	err = br.AddBlock(prevBlock)
	assert.NoError(t, err)

	eventService := common.NewEventService("", "Event")
	blockPool := blockchain.NewBlockPool()

	bApi, err := api.NewBlockApi(publisherID, br, eventService, blockPool)
	assert.NoError(t, err)

	commandHandler := adapter.NewBlockProposeCommandHandler(bApi, consensusService, "solo")

	//when
	_, errRPC := commandHandler.HandleProposeBlockCommand(command.ProposeBlock{TxList: nil})
	//then
	assert.Equal(t, errRPC, rpc.Error{Message: adapter.ErrCommandTransactions.Error()})

	//when
	_, errRPC = commandHandler.HandleProposeBlockCommand(command.ProposeBlock{TxList: make([]command.Tx, 0)})
	//then
	assert.Equal(t, errRPC, rpc.Error{Message: adapter.ErrCommandTransactions.Error()})

	//when
	_, errRPC = commandHandler.HandleProposeBlockCommand(command.ProposeBlock{
		TxList: []command.Tx{
			{
				ID:        "tx01",
				ICodeID:   "ICodeID",
				PeerID:    "2",
				TimeStamp: time.Now().Round(0),
				Jsonrpc:   "123",
				Function:  "function1",
				Args:      []string{"arg1", "arg2"},
				Signature: []byte{0x1},
			},
		},
	})

	//then
	assert.Equal(t, errRPC, rpc.Error{})

	//when
	_, errRPC = commandHandler.HandleProposeBlockCommand(command.ProposeBlock{
		TxList: []command.Tx{
			{
				ID:        "tx01",
				ICodeID:   "ICodeID",
				PeerID:    "2",
				TimeStamp: time.Now().Round(0),
				Jsonrpc:   "123",
				Function:  "function1",
				Args:      []string{"arg1", "arg2"},
				Signature: []byte{0x1},
			},

			{
				ID:        "tx02",
				ICodeID:   "ICodeID",
				PeerID:    "2",
				TimeStamp: time.Now().Round(0),
				Jsonrpc:   "123",
				Function:  "function1",
				Args:      []string{"arg1", "arg2"},
				Signature: []byte{0x1},
			},
		},
	})

	//then
	assert.Equal(t, errRPC, rpc.Error{})

	wg.Wait()
}

func TestBlockProposeCommandHandler_HandleProposeBlockCommand_NotSolo_OneTransaction(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	dbPath := "./.db"
	blockRepository, err := mem.NewBlockRepository(dbPath)
	assert.Equal(t, nil, err)
	defer func() {
		blockRepository.Close()
		os.RemoveAll(dbPath)
	}()

	prevBlock := mock.GetNewBlock([]byte("genesis"), 0)
	err = blockRepository.AddBlock(prevBlock)
	assert.NoError(t, err)

	consensusService := mock.ConsensusService{}
	consensusService.ConsensusBlockFunc = func(block blockchain.DefaultBlock) error {
		t.Log("consensus service")
		assert.Equal(t, uint64(1), block.GetHeight())
		assert.Equal(t, prevBlock.GetSeal(), block.GetPrevSeal())
		assert.Equal(t, []byte("iAmPublisher"), block.Creator)
		assert.Equal(t, 1, len(block.GetTxList()))
		wg.Done()
		return nil
	}

	publisherID := "iAmPublisher"
	eventService := common.NewEventService("", "Event")
	blockPool := blockchain.NewBlockPool()

	api, err := api.NewBlockApi(publisherID, blockRepository, eventService, blockPool)
	assert.NoError(t, err)

	commandHandler := adapter.NewBlockProposeCommandHandler(api, consensusService, "notSolo")

	//when
	_, errRPC := commandHandler.HandleProposeBlockCommand(command.ProposeBlock{
		TxList: []command.Tx{
			{
				ID:        "tx01",
				ICodeID:   "ICodeID",
				PeerID:    "2",
				TimeStamp: time.Now().Round(0),
				Jsonrpc:   "123",
				Function:  "function1",
				Args:      []string{"arg1", "arg2"},
				Signature: []byte{0x1},
			},
		},
	})
	//then
	assert.Equal(t, errRPC, rpc.Error{})

	wg.Wait()
}

func TestBlockProposeCommandHandler_HandleProposeBlockCommand_NotSolo_TwoTransaction(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)

	dbPath := "./.db"
	blockRepository, err := mem.NewBlockRepository(dbPath)
	assert.Equal(t, nil, err)
	defer func() {
		blockRepository.Close()
		os.RemoveAll(dbPath)
	}()

	prevBlock := mock.GetNewBlock([]byte("genesis"), 0)
	err = blockRepository.AddBlock(prevBlock)
	assert.NoError(t, err)

	consensusService := mock.ConsensusService{}
	consensusService.ConsensusBlockFunc = func(block blockchain.DefaultBlock) error {
		t.Log("consensus service")
		assert.Equal(t, uint64(1), block.GetHeight())
		assert.Equal(t, prevBlock.GetSeal(), block.GetPrevSeal())
		assert.Equal(t, []byte("iAmPublisher"), block.Creator)
		assert.Equal(t, 2, len(block.GetTxList()))
		wg.Done()
		return nil
	}

	publisherID := "iAmPublisher"
	eventService := common.NewEventService("", "Event")
	blockPool := blockchain.NewBlockPool()

	api, err := api.NewBlockApi(publisherID, blockRepository, eventService, blockPool)
	assert.NoError(t, err)

	commandHandler := adapter.NewBlockProposeCommandHandler(api, consensusService, "notSolo")

	//when
	_, errRPC := commandHandler.HandleProposeBlockCommand(command.ProposeBlock{
		TxList: []command.Tx{
			{
				ID:        "tx01",
				ICodeID:   "ICodeID",
				PeerID:    "2",
				TimeStamp: time.Now().Round(0),
				Jsonrpc:   "123",
				Function:  "function1",
				Args:      []string{"arg1", "arg2"},
				Signature: []byte{0x1},
			},
			{
				ID:        "tx02",
				ICodeID:   "ICodeID",
				PeerID:    "2",
				TimeStamp: time.Now().Round(0),
				Jsonrpc:   "123",
				Function:  "function1",
				Args:      []string{"arg1", "arg2"},
				Signature: []byte{0x1},
			},
		},
	})
	//then
	assert.Equal(t, errRPC, rpc.Error{})

	wg.Wait()
}

func TestBlockProposeCommandHandler_HandleProposeBlockCommand_NotSolo_ExceptionCases(t *testing.T) {
	ErrConsesnsusService := errors.New("Consensus Error")

	dbPath := "./.db"
	blockRepository, err := mem.NewBlockRepository(dbPath)
	assert.Equal(t, nil, err)
	defer func() {
		blockRepository.Close()
		os.RemoveAll(dbPath)
	}()

	prevBlock := mock.GetNewBlock([]byte("genesis"), 0)
	err = blockRepository.AddBlock(prevBlock)
	assert.NoError(t, err)

	consensusService := mock.ConsensusService{}
	consensusService.ConsensusBlockFunc = func(block blockchain.DefaultBlock) error {
		t.Log("consensus service")
		return ErrConsesnsusService
	}

	publisherID := "iAmPublisher"
	eventService := common.NewEventService("", "Event")
	blockPool := blockchain.NewBlockPool()

	api, err := api.NewBlockApi(publisherID, blockRepository, eventService, blockPool)
	assert.NoError(t, err)

	commandHandler := adapter.NewBlockProposeCommandHandler(api, consensusService, "notSolo")

	//when
	_, errRPC := commandHandler.HandleProposeBlockCommand(command.ProposeBlock{TxList: nil})
	//then
	assert.Equal(t, errRPC, rpc.Error{Message: adapter.ErrCommandTransactions.Error()})

	//when
	_, errRPC = commandHandler.HandleProposeBlockCommand(command.ProposeBlock{TxList: make([]command.Tx, 0)})
	//then
	assert.Equal(t, errRPC, rpc.Error{Message: adapter.ErrCommandTransactions.Error()})

	//when
	//when
	_, errRPC = commandHandler.HandleProposeBlockCommand(command.ProposeBlock{
		TxList: []command.Tx{
			{
				ID:        "tx01",
				ICodeID:   "ICodeID",
				PeerID:    "2",
				TimeStamp: time.Now().Round(0),
				Jsonrpc:   "123",
				Function:  "function1",
				Args:      []string{"arg1", "arg2"},
				Signature: []byte{0x1},
			},
		},
	})
	//then
	assert.Equal(t, errRPC, rpc.Error{Message: ErrConsesnsusService.Error()})

}
