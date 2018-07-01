package blockchain

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"time"

	"log"

	"github.com/stretchr/testify/assert"
)

var GenesisBlockConfigJson = []byte(`{
								  "Seal":[],
								  "PrevSeal":[],
								  "Height":0,
								  "TxList":[],
								  "TxSeal":[],
								  "TimeStamp":"0001-01-01T00:00:00-00:00",
								  "Creator":[]
								}`)

func TestCreateGenesisBlock(t *testing.T) {
	GenesisFilePath := "./GenesisBlockConfig.json"

	tests := map[string]struct {
		input  string
		output DefaultBlock
		err    error
	}{
		"success create genesisBlock": {
			input: GenesisFilePath,
			output: DefaultBlock{
				Height:    0,
				TxList:    make([]*DefaultTransaction, 0),
				TxSeal:    make([][]byte, 0),
				Timestamp: (time.Now()).Round(0),
				Creator:   make([]byte, 0),
			},
			err: nil,
		},
	}

	var tempBlock DefaultBlock
	err := json.Unmarshal(GenesisBlockConfigJson, &tempBlock)

	if err != nil {
		log.Println(err.Error())
	}

	GenesisBlockConfigByte, err := json.Marshal(tempBlock)

	err = ioutil.WriteFile(GenesisFilePath, GenesisBlockConfigByte, 0644)

	if err != nil {
		log.Println(err.Error())
	}

	defer os.Remove(GenesisFilePath)

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//when
		GenesisBlock, err := CreateGenesisBlock(GenesisFilePath)

		//then
		assert.Equal(t, test.err, err)
		assert.Equal(t, test.output.Height, GenesisBlock.Height)
		assert.Equal(t, test.output.TxList, GenesisBlock.TxList)
		assert.Equal(t, test.output.TxSeal, GenesisBlock.TxSeal)
		assert.Equal(t, test.output.Timestamp.String()[:19], GenesisBlock.Timestamp.String()[:19])
		assert.Equal(t, test.output.Creator, GenesisBlock.Creator)

	}

}
