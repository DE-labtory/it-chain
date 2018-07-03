package blockchain_test

import (
	"testing"

	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/stretchr/testify/assert"
)

func TestCreateGenesisBlock(t *testing.T) {

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
				TxList:    make([]blockchain.Transaction, 0),
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

			err: blockchain.ErrGetConfig,
		},
	}

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

	var tempBlock blockchain.DefaultBlock

	err := json.Unmarshal(GenesisBlockConfigJson, &tempBlock)

	if err != nil {
		log.Println(err.Error())
	}

	GenesisBlockConfigByte, err := json.Marshal(tempBlock)

	err = ioutil.WriteFile(GenesisFilePath, GenesisBlockConfigByte, 0644)

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
			Creator  []byte
		}
		output blockchain.Block
		err    error
	}{
		"success create proposed block": {

			input: struct {
				prevSeal []byte
				height   uint64
				txList   []blockchain.Transaction
				Creator  []byte
			}{
				prevSeal: []byte("prevseal"),
				height:   1,
				txList: []blockchain.Transaction{
					&blockchain.DefaultTransaction{},
				},
				Creator: []byte("junksound"),
			},

			output: &blockchain.DefaultBlock{
				PrevSeal: []byte("prevseal"),
				Height:   1,
				TxList: []blockchain.Transaction{
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
				Creator  []byte
			}{
				prevSeal: []byte("prevseal"),
				height:   1,
				txList:   []blockchain.Transaction{},
				Creator:  []byte("junksound"),
			},

			output: nil,

			err: blockchain.ErrBuildingTxSeal,
		},

		"fail case2: without prevseal or creator": {

			input: struct {
				prevSeal []byte
				height   uint64
				txList   []blockchain.Transaction
				Creator  []byte
			}{
				prevSeal: nil,
				height:   1,
				txList: []blockchain.Transaction{
					&blockchain.DefaultTransaction{},
				},
				Creator: nil,
			},

			output: nil,

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
			test.input.Creator,
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
