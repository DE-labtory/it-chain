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

	block := blockchain.CreateNewBlock(nil, "test")

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

	block := blockchain.CreateNewBlock(nil, "test")
	blockNumber := block.Header.Number

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

	block := blockchain.CreateNewBlock(nil, "test")
	blockHash := block.Header.BlockHash

	err := blockchainLevelDB.AddBlock(block)
	assert.NoError(t, err)

	retrievedBlock, err := blockchainLevelDB.GetBlockByHash(blockHash)
	assert.NoError(t, err)
	assert.Equal(t, block, retrievedBlock)
}