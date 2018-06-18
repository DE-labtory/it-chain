package blockchain

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/it-chain/yggdrasill/impl"
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
	wrongFilePath := "./WrongFileName.json"

	validator := new(impl.DefaultValidator)
	var tempBlock impl.DefaultBlock
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
	assert.Equal(t, make([]*impl.DefaultTransaction, 0), GenesisBlock.TxList)
	assert.Equal(t, make([][]byte, 0), GenesisBlock.TxSeal)
	assert.Equal(t, time.Now().String()[:19], GenesisBlock.Timestamp.String()[:19])
	assert.Equal(t, make([]byte, 0), GenesisBlock.Creator)

	_, err2 := CreateGenesisBlock(wrongFilePath)

	assert.Error(t, err2)
}
