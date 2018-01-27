package blockchaindb

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"it-chain/service/blockchain"
	"os"
)

func TestBlockchainDB_AddBlock(t *testing.T) {
	path := "./test_db"
	b := CreateNewBlockchainDB(path)
	defer func(){
		b.Close()
		os.RemoveAll(path)
	}()

	block := blockchain.CreateNewBlock(nil, "test")
	err := b.AddBlock(block)
	assert.NoError(t, err)
}

func TestBlockchainDBImpl_GetBlockByNumber(t *testing.T) {
	path := "./test_db"
	b := CreateNewBlockchainDB(path)
	defer func(){
		b.Close()
		os.RemoveAll(path)
	}()

	block := blockchain.CreateNewBlock(nil, "test")
	blockNumber := block.Header.Number

	err := b.AddBlock(block)
	assert.NoError(t, err)

	retrievedBlock, err := b.GetBlockByNumber(blockNumber)
	assert.NoError(t, err)
	assert.Equal(t, block, retrievedBlock)
}

func TestBlockchainDBImpl_GetBlockByHash(t *testing.T) {
	path := "./test_db"
	b := CreateNewBlockchainDB(path)
	defer func(){
		b.Close()
		os.RemoveAll(path)
	}()

	block := blockchain.CreateNewBlock(nil, "test")
	blockHash := block.Header.BlockHash

	err := b.AddBlock(block)
	assert.NoError(t, err)

	retrievedBlock, err := b.GetBlockByHash(blockHash)
	assert.NoError(t, err)
	assert.Equal(t, block, retrievedBlock)
}