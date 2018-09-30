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

	"sync"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/engine/txpool/api"
	"github.com/it-chain/engine/txpool/infra/mem"
	"github.com/it-chain/engine/txpool/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestTransactionApi_CreateTransaction(t *testing.T) {

	tests := map[string]struct {
		input struct {
			txData txpool.TxData
		}
		err error
	}{
		"success": {
			input: struct {
				txData txpool.TxData
			}{
				txData: txpool.TxData{
					ICodeID:   "gg",
					Function:  "1",
					Signature: []byte("123"),
					Args:      []string{"1", "2"},
					Jsonrpc:   "2.0",
				},
			},
			err: nil,
		},
	}

	transactionRepository := mem.NewTransactionRepository()
	leaderRepository := mem.NewLeaderRepository()
	eventService := mock.EventService{}
	transferService := txpool.NewTransferService(transactionRepository, leaderRepository, eventService)
	blockProposalService := txpool.NewBlockProposalService(transactionRepository, eventService)
	transactionApi := api.NewTransactionApi("zf", transactionRepository, leaderRepository, transferService, blockProposalService)

	for _, test := range tests {
		tx, err := transactionApi.CreateTransaction(test.input.txData)

		assert.Equal(t, test.err, err)
		assert.Equal(t, tx.ICodeID, test.input.txData.ICodeID)
		assert.Equal(t, tx.Args, test.input.txData.Args)
		assert.Equal(t, tx.Signature, test.input.txData.Signature)
		assert.Equal(t, tx.Jsonrpc, test.input.txData.Jsonrpc)
		assert.Equal(t, tx.Function, test.input.txData.Function)
	}
}

func TestTransactionApi_DeleteTransaction(t *testing.T) {

	tests := map[string]struct {
		input string
		err   error
	}{
		"success": {
			input: "transactionID",
			err:   mem.ErrTransactionDoesNotExist,
		},
	}

	transactionRepository := mem.NewTransactionRepository()
	leaderRepository := mem.NewLeaderRepository()
	eventService := mock.EventService{}
	transferService := txpool.NewTransferService(transactionRepository, leaderRepository, eventService)
	blockProposalService := txpool.NewBlockProposalService(transactionRepository, eventService)
	transactionApi := api.NewTransactionApi("zf", transactionRepository, leaderRepository, transferService, blockProposalService)

	transactionRepository.Save(txpool.Transaction{
		ID: "transactionID",
	})

	for _, test := range tests {
		transactionApi.DeleteTransaction(test.input)

		_, err := transactionRepository.FindById(test.input)
		assert.Equal(t, err, test.err)
	}
}

func TestTransactionApi_ProposeBlock_Solo(t *testing.T) {

	//publish 하는 걸 잘 받는지.

	//set api
	//mode가 솔로일 때, pbft일 때
	wg := sync.WaitGroup{}
	wg.Add(1)

	tests := map[string]struct {
		engineMode string
		txList     []txpool.Transaction
		wgNum      int
	}{
		"success-solo": {
			engineMode: "solo",
			txList: []txpool.Transaction{
				{
					ID: "tx01",
				},
				{
					ID: "tx02",
				},
			},
		},
		//"success-no transaction": {
		//	engineMode: "solo",
		//	txList:     []txpool.Transaction{},
		//},
	}

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		handler := &mock.ProposeEventHandler{}
		handler.HandleFunc = func(command command.ProposeBlock) {
			assert.Equal(t, len(test.txList), len(command.TxList))
			wg.Done()
		}

		subscriber.SubscribeTopic("block.*", handler)

		//set repo
		txPoolRepo := mem.NewTransactionRepository()

		for _, tx := range test.txList {
			txPoolRepo.Save(tx)
		}

		leaderRepo := mem.NewLeaderRepository()
		eventService := common.NewEventService("", "Event")

		//set service
		transferService := txpool.NewTransferService(txPoolRepo, leaderRepo, eventService)
		blockProposalService := txpool.NewBlockProposalService(txPoolRepo, eventService)

		//set api
		transactionApi := api.NewTransactionApi("node01", txPoolRepo, leaderRepo, transferService, blockProposalService)

		err := transactionApi.ProposeBlock(test.engineMode)

		assert.NoError(t, err)
	}

	wg.Wait()

}

func TestTransactionApi_ProposeBlock_Solo_NoTransaction(t *testing.T) {

	tests := map[string]struct {
		engineMode string
		txList     []txpool.Transaction
		wgNum      int
	}{
		"success-no transaction": {
			engineMode: "solo",
			txList:     []txpool.Transaction{},
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//set repo
		txPoolRepo := mem.NewTransactionRepository()

		for _, tx := range test.txList {
			txPoolRepo.Save(tx)
		}

		leaderRepo := mem.NewLeaderRepository()
		eventService := common.NewEventService("", "Event")

		//set service
		transferService := txpool.NewTransferService(txPoolRepo, leaderRepo, eventService)
		blockProposalService := txpool.NewBlockProposalService(txPoolRepo, eventService)

		//set api
		transactionApi := api.NewTransactionApi("node01", txPoolRepo, leaderRepo, transferService, blockProposalService)

		err := transactionApi.ProposeBlock(test.engineMode)

		assert.NoError(t, err)
	}

}

func TestTransactionApi_ProposeBlock_PBFT(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	tests := map[string]struct {
		engineMode string
		txList     []txpool.Transaction
		leader     txpool.Leader
	}{
		"success-pbft": {
			engineMode: "pbft",
			txList: []txpool.Transaction{
				{
					ID: "tx03",
				},
				{
					ID: "tx04",
				},
				{
					ID: "tx05",
				},
			},

			leader: txpool.Leader{Id: "leader"},
		},
	}

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		handler := &mock.ProposeEventHandler{}
		handler.HandleFunc = func(command command.ProposeBlock) {
			assert.Equal(t, len(test.txList), len(command.TxList))
			wg.Done()
		}

		subscriber.SubscribeTopic("block.*", handler)

		//set repo
		txPoolRepo := mem.NewTransactionRepository()

		for _, tx := range test.txList {
			txPoolRepo.Save(tx)
		}

		leaderRepo := mem.NewLeaderRepository()
		leaderRepo.Set(test.leader)
		eventService := common.NewEventService("", "Event")

		//set service
		transferService := txpool.NewTransferService(txPoolRepo, leaderRepo, eventService)
		blockProposalService := txpool.NewBlockProposalService(txPoolRepo, eventService)

		//set api
		transactionApi := api.NewTransactionApi("leader", txPoolRepo, leaderRepo, transferService, blockProposalService)

		err := transactionApi.ProposeBlock(test.engineMode)

		assert.NoError(t, err)
	}

	wg.Wait()

}

func TestTransactionApi_SendLeaderTransaction(t *testing.T) {

	wg := sync.WaitGroup{}
	wg.Add(1)

	tests := map[string]struct {
		engineMode string
		txList     []txpool.Transaction
		leader     txpool.Leader
	}{
		"success-pbft": {
			engineMode: "pbft",
			txList: []txpool.Transaction{
				{
					ID: "tx03",
				},
				{
					ID: "tx04",
				},
				{
					ID: "tx05",
				},
			},

			leader: txpool.Leader{Id: "leader"},
		},
	}

	subscriber := pubsub.NewTopicSubscriber("", "Event")
	defer subscriber.Close()

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		handler := &mock.SendTransactionCommandHandler{}
		handler.HandleFunc = func(command command.DeliverGrpc) {
			assert.Equal(t, "SendLeaderTransactionsProtocol", command.Protocol)
			assert.Equal(t, test.leader.Id, command.RecipientList[0])
			wg.Done()

		}

		subscriber.SubscribeTopic("message.*", handler)

		//set repo
		txPoolRepo := mem.NewTransactionRepository()

		for _, tx := range test.txList {
			txPoolRepo.Save(tx)

		}

		leaderRepo := mem.NewLeaderRepository()
		leaderRepo.Set(test.leader)
		eventService := common.NewEventService("", "Event")

		//set service
		transferService := txpool.NewTransferService(txPoolRepo, leaderRepo, eventService)
		blockProposalService := txpool.NewBlockProposalService(txPoolRepo, eventService)

		//set api
		transactionApi := api.NewTransactionApi("node01", txPoolRepo, leaderRepo, transferService, blockProposalService)

		err := transactionApi.SendLeaderTransaction(test.engineMode)
		assert.NoError(t, err)
	}

	wg.Wait()
}
