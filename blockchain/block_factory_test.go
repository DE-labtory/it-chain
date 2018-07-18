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

package blockchain_test

import (
	"testing"

	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/test/mock"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

func TestCreateGenesisBlock(t *testing.T) {

	//given

	const shortForm = "2006-Jan-02"
	timeStamp, _ := time.Parse(shortForm, "0000-Jan-00")

	tests := map[string]struct {
		input struct {
			ConfigFilePath string
		}
		output blockchain.Block
		err    error
	}{
		"success create genesisBlock": {

			input: struct {
				ConfigFilePath string
			}{
				ConfigFilePath: "./GenesisBlockConfig.json",
			},

			output: &blockchain.DefaultBlock{
				PrevSeal:  make([]byte, 0),
				Height:    uint64(0),
				TxList:    make([]*blockchain.DefaultTransaction, 0),
				TxSeal:    make([][]byte, 0),
				Timestamp: timeStamp,
				Creator:   make([]byte, 0),
			},

			err: nil,
		},

		"fail create genesisBlock: wrong file path": {

			input: struct {
				ConfigFilePath string
			}{
				ConfigFilePath: "./WrongBlockConfig.json",
			},

			output: nil,

			err: blockchain.ErrSetConfig,
		},
	}

	repo := mock.EventRepository{}

	repo.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, 1, len(events))
		assert.IsType(t, &blockchain.BlockCreatedEvent{}, events[0])
		return nil
	}
	repo.CloseFunc = func() {}

	eventstore.InitForMock(repo)
	defer eventstore.Close()

	GenesisFilePath := "./GenesisBlockConfig.json"

	defer os.Remove(GenesisFilePath)

	GenesisBlockConfigJson := []byte(`{
								  "Seal":[],
								  "PrevSeal":[],
								  "Height":0,
								  "TxList":[],
								  "TxSeal":[],
								  "TimeStamp":"0001-01-01T00:00:00-00:00",
								  "Creator":[]
								}`)

	err := ioutil.WriteFile(GenesisFilePath, GenesisBlockConfigJson, 0644)

	if err != nil {
		log.Println(err.Error())
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//when
		GenesisBlock, err := blockchain.CreateGenesisBlock(test.input.ConfigFilePath)

		//then
		assert.Equal(t, test.err, err)

		if err != nil {
			assert.Equal(t, test.output, GenesisBlock)
			continue
		}

		assert.Equal(t, test.output.GetPrevSeal(), GenesisBlock.GetPrevSeal())
		assert.Equal(t, test.output.GetHeight(), GenesisBlock.GetHeight())
		assert.Equal(t, test.output.GetTxList(), GenesisBlock.GetTxList())
		assert.Equal(t, test.output.GetTxSeal(), GenesisBlock.GetTxSeal())
		assert.Equal(t, test.output.GetTimestamp().String()[:19], GenesisBlock.GetTimestamp().String()[:19])
		assert.Equal(t, test.output.GetCreator(), GenesisBlock.GetCreator())

	}

}

func TestCreateProposedBlock(t *testing.T) {

	//given

	timeStamp := time.Now().Round(0)

	tests := map[string]struct {
		input struct {
			prevSeal []byte
			height   uint64
			txList   []blockchain.Transaction
			creator  []byte
		}
		output blockchain.Block
		err    error
	}{
		"success create proposed block": {

			input: struct {
				prevSeal []byte
				height   uint64
				txList   []blockchain.Transaction
				creator  []byte
			}{
				prevSeal: []byte("prevseal"),
				height:   1,
				txList: []blockchain.Transaction{
					&blockchain.DefaultTransaction{
						ID:        "tx01",
						Status:    0,
						PeerID:    "junksound",
						Timestamp: timeStamp,
						TxData: &blockchain.TxData{
							Jsonrpc: "",
							Method:  "",
							Params:  blockchain.Params{},
							ID:      "txData01",
						},
						Signature: []byte("Signature"),
					},
				},
				creator: []byte("junksound"),
			},

			output: &blockchain.DefaultBlock{
				PrevSeal: []byte("prevseal"),
				Height:   1,
				TxList: []*blockchain.DefaultTransaction{
					{
						ID:        "tx01",
						Status:    0,
						PeerID:    "junksound",
						Timestamp: timeStamp,
						TxData: &blockchain.TxData{
							Jsonrpc: "",
							Method:  "",
							Params:  blockchain.Params{},
							ID:      "txData01",
						},
						Signature: []byte("Signature"),
					},
				},
				Timestamp: timeStamp,
				Creator:   []byte("junksound"),
			},

			err: nil,
		},

		"fail case1: without transaction": {

			input: struct {
				prevSeal []byte
				height   uint64
				txList   []blockchain.Transaction
				creator  []byte
			}{
				prevSeal: []byte("prevseal"),
				height:   1,
				txList:   nil,
				creator:  []byte("junksound"),
			},

			output: nil,

			err: blockchain.ErrBuildingTxSeal,
		},

		"fail case2: without prevseal or creator": {

			input: struct {
				prevSeal []byte
				height   uint64
				txList   []blockchain.Transaction
				creator  []byte
			}{
				prevSeal: nil,
				height:   1,
				txList: []blockchain.Transaction{
					&blockchain.DefaultTransaction{
						ID:        "tx01",
						Status:    0,
						PeerID:    "junksound",
						Timestamp: timeStamp,
						TxData: &blockchain.TxData{
							Jsonrpc: "",
							Method:  "",
							Params:  blockchain.Params{},
							ID:      "txData01",
						},
						Signature: []byte("Signature"),
					},
				},
				creator: nil,
			},

			output: nil,

			err: blockchain.ErrBuildingSeal,
		},
	}

	repo := mock.EventRepository{}

	repo.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, 1, len(events))
		assert.IsType(t, &blockchain.BlockCreatedEvent{}, events[0])
		return nil
	}
	repo.CloseFunc = func() {}

	eventstore.InitForMock(repo)
	defer eventstore.Close()

	for testName, test := range tests {

		t.Logf("Running test case %s", testName)

		//when
		ProposedBlock, err := blockchain.CreateProposedBlock(
			test.input.prevSeal,
			test.input.height,
			test.input.txList,
			test.input.creator,
		)

		//then
		assert.Equal(t, test.err, err)

		if err != nil {
			assert.Equal(t, test.output, ProposedBlock)
			continue
		}

		assert.Equal(t, test.output.GetPrevSeal(), ProposedBlock.GetPrevSeal())
		assert.Equal(t, test.output.GetHeight(), ProposedBlock.GetHeight())
		assert.Equal(t, test.output.GetTxList(), ProposedBlock.GetTxList())
		assert.Equal(t, test.output.GetTimestamp().String()[:19], ProposedBlock.GetTimestamp().String()[:19])
		assert.Equal(t, test.output.GetCreator(), ProposedBlock.GetCreator())
	}

}

func TestCreateRetrievedBlock(t *testing.T) {

	//given
	timeStamp := time.Now().Round(0)
	prevSeal := []byte("prevseal")
	height := uint64(0)
	txList := []blockchain.Transaction{
		&blockchain.DefaultTransaction{
			ID:        "tx01",
			Status:    0,
			PeerID:    "junksound",
			Timestamp: timeStamp,
			TxData: &blockchain.TxData{
				Jsonrpc: "",
				Method:  "",
				Params:  blockchain.Params{},
				ID:      "txData01",
			},
			Signature: []byte("Signature"),
		},
	}
	creator := []byte("junksound")

	retrievedBlock, err := blockchain.CreateProposedBlock(prevSeal, height, txList, creator)
	if err != nil {
	}

	tests := map[string]struct {
		input struct {
			retrivedBlock blockchain.Block
		}
		output struct {
			createdBlock blockchain.Block
		}
		err error
	}{
		"success create retrieved block": {
			input: struct {
				retrivedBlock blockchain.Block
			}{
				retrivedBlock: retrievedBlock,
			},

			output: struct {
				createdBlock blockchain.Block
			}{
				createdBlock: retrievedBlock,
			},

			err: nil,
		},
	}

	repo := mock.EventRepository{}

	repo.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, 1, len(events))
		assert.IsType(t, &blockchain.BlockCreatedEvent{}, events[0])
		return nil
	}
	repo.CloseFunc = func() {}

	eventstore.InitForMock(repo)
	defer eventstore.Close()

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//when
		RetrivedBlock, err := blockchain.CreateRetrievedBlock(test.input.retrivedBlock)
		assert.Equal(t, test.err, err)
		assert.Equal(t, test.output.createdBlock, RetrivedBlock)

	}
}
