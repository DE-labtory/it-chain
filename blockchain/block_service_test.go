package blockchain

import (
	"go/build"
	"testing"

	"encoding/json"

	"io/ioutil"

	"os"

	"time"

	"github.com/it-chain/yggdrasill/impl"
	"github.com/stretchr/testify/assert"
)

func TestCreateGenesisBlock(t *testing.T) {
	genesisconfPath := build.Default.GOPATH + "/src/github.com/it-chain/it-chain-Engine/.it-chain/genesisconf/"
	genesisConfFilePath := genesisconfPath + "GenesisBlockConfig.json"
	tempFilePath := genesisconfPath + "TempBlockConfig.json"
	wrongFilePath := genesisconfPath + "WrongFileName.json"
	tempBlockConfigJson := []byte(`{
								  "Seal":[],
								  "PrevSeal":[],
								  "Height":0,
								  "TxList":[],
								  "TxSeal":[],
								  "TimeStamp":"0001-01-01T00:00:00-00:00",
								  "Creator":[]
								}`)
	validator := new(impl.DefaultValidator)
	var tempBlock impl.DefaultBlock
	_ = json.Unmarshal(tempBlockConfigJson, &tempBlock)
	tempBlockConfigByte, _ := json.Marshal(tempBlock)
	_ = ioutil.WriteFile(tempFilePath, tempBlockConfigByte, 0644)
	defer os.Remove(tempFilePath)
	rightFilePaths := []string{genesisConfFilePath, tempFilePath}
	for _, rightFilePath := range rightFilePaths {
		GenesisBlock, err1 := CreateGenesisBlock(rightFilePath)
		expectedSeal, _ := validator.BuildSeal(GenesisBlock)
		assert.NoError(t, err1)
		assert.Equal(t, expectedSeal, GenesisBlock.Seal)
		assert.Equal(t, make([]byte, 0), GenesisBlock.PrevSeal)
		assert.Equal(t, uint64(0), GenesisBlock.Height)
		assert.Equal(t, make([]*impl.DefaultTransaction, 0), GenesisBlock.TxList)
		assert.Equal(t, make([][]byte, 0), GenesisBlock.TxSeal)
		assert.Equal(t, time.Now().String()[:19], GenesisBlock.Timestamp.String()[:19])
		assert.Equal(t, make([]byte, 0), GenesisBlock.Creator)
	}
	_, err2 := CreateGenesisBlock(wrongFilePath)
	assert.Error(t, err2)
}
