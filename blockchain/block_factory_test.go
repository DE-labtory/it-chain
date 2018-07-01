package blockchain_test

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"log"

	"time"

	"os"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/stretchr/testify/assert"
)

func TestCreateGenesisBlock(t *testing.T) {

	tests := map[string]struct {
		input struct {
			ConfigFilePath string
		}
		output *blockchain.DefaultBlock
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

		assert.Equal(t, test.output.Height, GenesisBlock.GetHeight())
		assert.Equal(t, test.output.TxList, GenesisBlock.GetTxList())
		assert.Equal(t, test.output.TxSeal, GenesisBlock.GetTxSeal())
		assert.Equal(t, test.output.Timestamp.String()[:19], GenesisBlock.GetTimestamp().String()[:19])
		assert.Equal(t, test.output.Creator, GenesisBlock.GetCreator())

	}

}
