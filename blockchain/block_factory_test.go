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

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

type MockRepostiory struct {
	loadFunc func(aggregate midgard.Aggregate, aggregateID string) error
	saveFunc func(aggregateID string, events ...midgard.Event) error
}

func (m MockRepostiory) Load(aggregate midgard.Aggregate, aggregateID string) error {
	return m.loadFunc(aggregate, aggregateID)
}

func (m MockRepostiory) Save(aggregateID string, events ...midgard.Event) error {
	return m.saveFunc(aggregateID, events...)
}

func (MockRepostiory) Close() {}

func TestCreateGenesisBlock(t *testing.T) {
	//given

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
				Timestamp: (time.Now()).Round(0),
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

	repo := MockRepostiory{}

	repo.saveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, 1, len(events))
		assert.IsType(t, &blockchain.BlockCreatedEvent{}, events[0])
		return nil
	}

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
					&blockchain.DefaultTransaction{},
				},
				creator: []byte("junksound"),
			},

			output: &blockchain.DefaultBlock{
				PrevSeal: []byte("prevseal"),
				Height:   1,
				TxList: []*blockchain.DefaultTransaction{
					&blockchain.DefaultTransaction{},
				},
				Timestamp: (time.Now()).Round(0),
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
					&blockchain.DefaultTransaction{},
				},
				creator: nil,
			},

			output: nil,

			err: blockchain.ErrBuildingSeal,
		},
	}

	repo := MockRepostiory{}

	repo.saveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, 1, len(events))
		assert.IsType(t, &blockchain.BlockCreatedEvent{}, events[0])
		return nil
	}

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
