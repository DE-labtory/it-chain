package blockchainleveldb

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"it-chain/service/blockchain"
)

func TestBlockchainLevelDB_AddBlock(t *testing.T) {
	path := "./test_db"
	blockchainLevelDB := CreateNewBlockchainLevelDB(path)
	defer func(){
		blockchainLevelDB.Close()
		os.RemoveAll(path)
	}()

	block := &blockchain.Block{}

	err := blockchainLevelDB.AddBlock(block)
	assert.NoError(t, err)
}

func TestBlockchainLevelDB_GetBlockByNumber(t *testing.T) {
	path := "./test_db"
	blockchainLevelDB := CreateNewBlockchainLevelDB(path)
	defer func(){
		blockchainLevelDB.Close()
		os.RemoveAll(path)
	}()

	block := &blockchain.Block{}

	blockNumber := uint64(1)
	block.Header.Number = blockNumber

	err := blockchainLevelDB.AddBlock(block)
	assert.NoError(t, err)

	retrievedBlock, err := blockchainLevelDB.GetBlockByNumber(blockNumber)
	assert.NoError(t, err)
	assert.Equal(t, block, retrievedBlock)
}

func TestBlockchainLevelDB_GetBlockByHash(t *testing.T) {
	path := "./test_db"
	blockchainLevelDB := CreateNewBlockchainLevelDB(path)
	defer func(){
		blockchainLevelDB.Close()
		os.RemoveAll(path)
	}()

	block := &blockchain.Block{}

	blockHash := "hash"
	block.Header.DataHash = blockHash

	err := blockchainLevelDB.AddBlock(block)
	assert.NoError(t, err)

	retrievedBlock, err := blockchainLevelDB.GetBlockByHash(blockHash)
	assert.NoError(t, err)
	assert.Equal(t, block, retrievedBlock)
}