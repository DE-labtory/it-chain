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
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/it-chain/engine/blockchain"
	"github.com/stretchr/testify/assert"
)

func TestCreateGenesisBlock(t *testing.T) {

	//given
	const longForm = "Jan 1, 2006 at 0:00am (MST)"
	timeStamp, _ := time.Parse(longForm, "Jan 1, 2018 at 0:00am (KST)")
	tree := &blockchain.DefaultTree{}
	tree.SetTxSealRoot([]byte("genesis"))

	tests := map[string]struct {
		input struct {
			ConfigFilePath string
		}
		output blockchain.DefaultBlock
		err    error
	}{
		"success create genesisBlock": {

			input: struct {
				ConfigFilePath string
			}{
				ConfigFilePath: "./GenesisBlockConfig.json",
			},

			output: blockchain.DefaultBlock{
				PrevSeal:  make([]byte, 0),
				Height:    uint64(0),
				TxList:    make([]*blockchain.DefaultTransaction, 0),
				Tree:      tree,
				Timestamp: timeStamp,
				Creator:   "junksound",
				State:     blockchain.Created,
			},

			err: nil,
		},

		"fail create genesisBlock: wrong file path": {

			input: struct {
				ConfigFilePath string
			}{
				ConfigFilePath: "./WrongBlockConfig.json",
			},

			output: blockchain.DefaultBlock{},

			err: blockchain.ErrSetConfig,
		},
	}

	GenesisFilePath := "./GenesisBlockConfig.json"

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
		assert.Equal(t, test.output.GetTimestamp().String()[:19], GenesisBlock.GetTimestamp().String()[:19])
		assert.Equal(t, test.output.GetCreator(), GenesisBlock.GetCreator())
		assert.Equal(t, test.output.GetState(), GenesisBlock.GetState())
		assert.Equal(t, test.output.Tree, GenesisBlock.Tree)
	}

}

func TestCreateProposedBlock(t *testing.T) {

	//given

	timeStamp := time.Now().Round(0)

	tests := map[string]struct {
		input struct {
			prevSeal []byte
			height   uint64
			txList   []*blockchain.DefaultTransaction
			creator  string
		}
		output blockchain.DefaultBlock
		err    error
	}{
		"success create proposed block": {

			input: struct {
				prevSeal []byte
				height   uint64
				txList   []*blockchain.DefaultTransaction
				creator  string
			}{
				prevSeal: []byte("prevseal"),
				height:   1,
				txList: []*blockchain.DefaultTransaction{
					{
						ID:        "tx01",
						ICodeID:   "ICodeID",
						PeerID:    "junksound",
						Timestamp: timeStamp,
						Jsonrpc:   "",
						Function:  "",
						Args:      make([]string, 0),
						Signature: []byte("Signature"),
					},
				},
				creator: "junksound",
			},

			output: blockchain.DefaultBlock{
				PrevSeal: []byte("prevseal"),
				Height:   1,
				TxList: []*blockchain.DefaultTransaction{
					{
						ID:        "tx01",
						ICodeID:   "ICodeID",
						PeerID:    "junksound",
						Timestamp: timeStamp,
						Jsonrpc:   "",
						Function:  "",
						Args:      make([]string, 0),
						Signature: []byte("Signature"),
					},
				},
				Timestamp: timeStamp,
				Creator:   "junksound",
				State:     blockchain.Created,
			},

			err: nil,
		},

		"fail case1: without transaction": {

			input: struct {
				prevSeal []byte
				height   uint64
				txList   []*blockchain.DefaultTransaction
				creator  string
			}{
				prevSeal: []byte("prevseal"),
				height:   1,
				txList:   nil,
				creator:  "junksound",
			},

			output: blockchain.DefaultBlock{},

			err: blockchain.ErrEmptyTxList,
		},

		"fail case2: without prevseal or creator": {

			input: struct {
				prevSeal []byte
				height   uint64
				txList   []*blockchain.DefaultTransaction
				creator  string
			}{
				prevSeal: nil,
				height:   1,
				txList: []*blockchain.DefaultTransaction{
					{
						ID:        "tx01",
						ICodeID:   "ICodeID",
						PeerID:    "junksound",
						Timestamp: timeStamp,
						Jsonrpc:   "",
						Function:  "",
						Args:      make([]string, 0),
						Signature: []byte("Signature"),
					},
				},
				creator: "",
			},

			output: blockchain.DefaultBlock{},

			err: blockchain.ErrBuildingSeal,
		},
	}

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
		assert.Equal(t, test.output.GetTxList()[0].GetID(), ProposedBlock.GetTxList()[0].GetID())
		assert.Equal(t, test.output.GetTimestamp().String()[:19], ProposedBlock.GetTimestamp().String()[:19])
		assert.Equal(t, test.output.GetCreator(), ProposedBlock.GetCreator())
		assert.Equal(t, test.output.GetState(), ProposedBlock.GetState())

	}

}
