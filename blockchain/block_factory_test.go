package blockchain

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"time"

	"fmt"

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

	tests := map[string]struct {
		input  string
		output DefaultBlock
		err    error
	}{
		"success:": {
			input: "./GenesisBlockConfig.json",
			output: DefaultBlock{
				Height:    0,
				TxList:    nil,
				TxSeal:    nil,
				Timestamp: time.Now(),
				Creator:   []byte("ddd"),
			},
			err: nil,
		},
		"fail:": {
			input: "./GenesisBlockConfig.json",
			output: DefaultBlock{
				Height:    0,
				TxList:    nil,
				TxSeal:    nil,
				Timestamp: time.Now(),
				Creator:   nil,
			},
			err: nil,
		},
	}
	fmt.Println("tests success")
	if tests["success:"].output.TxList == nil {
		fmt.Println("kkkkk")
	}

	GenesisFilePath := "./GenesisBlockConfig.json"
	wrongFilePath := "./WrongFileName.json"

	validator := new(DefaultValidator)
	var tempBlock DefaultBlock
	err := json.Unmarshal(GenesisBlockConfigJson, &tempBlock)
	assert.NoError(t, err)

	GenesisBlockConfigByte, _ := json.Marshal(tempBlock)
	err = ioutil.WriteFile(GenesisFilePath, GenesisBlockConfigByte, 0644)
	assert.NoError(t, err)

	defer os.Remove(GenesisFilePath)

	GenesisBlock, err1 := CreateGenesisBlock(GenesisFilePath)
	expectedSeal, _ := validator.BuildSeal(GenesisBlock)
	assert.NoError(t, err1)
	assert.Equal(t, expectedSeal, GenesisBlock.Seal)
	assert.Equal(t, make([]byte, 0), GenesisBlock.PrevSeal)
	assert.Equal(t, uint64(0), GenesisBlock.Height)
	assert.Equal(t, make([]*DefaultTransaction, 0), GenesisBlock.TxList)
	assert.Equal(t, make([][]byte, 0), GenesisBlock.TxSeal)
	assert.Equal(t, time.Now().String()[:19], GenesisBlock.Timestamp.String()[:19])
	assert.Equal(t, make([]byte, 0), GenesisBlock.Creator)
	fmt.Println("GenesisBlock")
	if GenesisBlock.TxList == nil {
		fmt.Println("tttttt")
	}

	_, err2 := CreateGenesisBlock(wrongFilePath)

	assert.Error(t, err2)
}

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
