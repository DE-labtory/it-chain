package db

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"os"
	"github.com/it-chain/it-chain-Engine/blockchain/domain/model"
)

func TestBlockchainDB_AddBlock(t *testing.T) {
	path := "./test_db"
	b := CreateNewBlockchainDB(path)
	defer func(){
		b.Close()
		os.RemoveAll(path)
	}()

	block := model.CreateNewBlock(nil)
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

	block := model.CreateNewBlock(nil)
	blockNumber := block.BlockHeader.Number

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

	block := model.CreateNewBlock(nil)
	blockHash := block.BlockHeader.BlockHash

	err := b.AddBlock(block)
	assert.NoError(t, err)

	retrievedBlock, err := b.GetBlockByHash(blockHash)
	assert.NoError(t, err)
	assert.Equal(t, block, retrievedBlock)

}