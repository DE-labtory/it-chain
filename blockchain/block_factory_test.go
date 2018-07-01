package blockchain

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"time"

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
		"success:": {
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

	//validator := new(DefaultValidator)
	var tempBlock DefaultBlock
	err := json.Unmarshal(GenesisBlockConfigJson, &tempBlock)
	assert.NoError(t, err)

	GenesisBlockConfigByte, _ := json.Marshal(tempBlock)
	err = ioutil.WriteFile(GenesisFilePath, GenesisBlockConfigByte, 0644)
	assert.NoError(t, err)

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

// fail case

func TestConfigFromJson(t *testing.T) {

	genesisFilePath := "./GenesisBlockConfig.json"
	wrongFilePath := "./WrongFileName.json"
	var tempJson DefaultBlock

	err := json.Unmarshal(GenesisBlockConfigJson, &tempJson)
	assert.NoError(t, err)

	GenesisBlockConfigByte, _ := json.Marshal(tempJson)
	err = ioutil.WriteFile(genesisFilePath, GenesisBlockConfigByte, 0644)
	assert.NoError(t, err)

	defer os.Remove(genesisFilePath)

	_, err1 := ConfigFromJson(genesisFilePath)
	assert.NoError(t, err1)
	_, err2 := ConfigFromJson(wrongFilePath)
	assert.Error(t, err2)
}
