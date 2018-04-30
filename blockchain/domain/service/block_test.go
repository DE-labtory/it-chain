package service

import (
	"go/build"
	"testing"
	"time"

	"encoding/json"

	"io/ioutil"

	"os"

	"github.com/it-chain/yggdrasill/block"
	tx "github.com/it-chain/yggdrasill/transaction"
	"github.com/stretchr/testify/assert"
)

func TestCreateGenesisBlock(t *testing.T) {
	genesisconfPath := build.Default.GOPATH + "/src/github.com/it-chain/it-chain-Engine/.it-chain/genesisconf/"
	rightFilePath := genesisconfPath + "TempBlockConfig.json"
	wrongFilePath := genesisconfPath + "WrongFileName.json"
	tempBlockConfigJson := []byte(`{"Header": {
										  "Height":0,
										  "PreviousHash":"",
										  "Version":"",
										  "MerkleTreeRootHash":"",
										  "TimeStamp":"0001-01-01T00:00:00-00:00",
										  "CreatorID":"Genesis",
										  "Signature":[],
										  "BlockHash":"",
										  "MerkleTreeHeight":0,
										  "TransactionCount":0
										},
							  "Proof": [],
							  "Transactions":[]
							}`)

	var tempBlock block.DefaultBlock
	_ = json.Unmarshal(tempBlockConfigJson, &tempBlock)
	tempBlockConfigByte, _ := json.Marshal(tempBlock)
	_ = ioutil.WriteFile(rightFilePath, tempBlockConfigByte, 0644)
	defer os.Remove(rightFilePath)

	GenesisBlock, err1 := CreateGenesisBlock(rightFilePath)
	_, err2 := CreateGenesisBlock(wrongFilePath)
	assert.NoError(t, err1)
	assert.Error(t, err2)
	assert.Equal(t, uint64(0), GenesisBlock.Header.Height)
	assert.Equal(t, "", GenesisBlock.Header.PreviousHash)
	assert.Equal(t, "", GenesisBlock.Header.Version)
	assert.Equal(t, "", GenesisBlock.Header.MerkleTreeRootHash)
	assert.Equal(t, time.Now().String()[:19], GenesisBlock.Header.TimeStamp.String()[:19])
	assert.Equal(t, "Genesis", GenesisBlock.Header.CreatorID)
	assert.Equal(t, make([]byte, 0), GenesisBlock.Header.Signature)
	assert.Equal(t, "", GenesisBlock.Header.BlockHash)
	assert.Equal(t, 0, GenesisBlock.Header.MerkleTreeHeight)
	assert.Equal(t, 0, GenesisBlock.Header.TransactionCount)
	assert.Equal(t, make([][]byte, 0), GenesisBlock.Proof)
	assert.Equal(t, make([]*tx.DefaultTransaction, 0), GenesisBlock.Transactions)
}
