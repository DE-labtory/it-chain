package repository

import (
	"go/build"
	"testing"

	"os"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/yggdrasill/impl"
	"github.com/stretchr/testify/assert"
)

var genesisConfFilePath = build.Default.GOPATH + "/src/github.com/it-chain/it-chain-Engine/.it-chain/genesisconf/GenesisBlockConfig.json"

func TestBlockRepository_AddBlock(t *testing.T) {
	dbPath := "./.db"
	db := leveldbwrapper.CreateNewDB(dbPath)
	defer func() {
		db.Close()
		os.RemoveAll(dbPath)
	}()

	genesisBlock, _ := blockchain.CreateGenesisBlock(genesisConfFilePath)
	validator := new(impl.DefaultValidator)
	opts := map[string]interface{}{
		"db_path": dbPath,
	}
	br := NewBlockRepository(db, *validator, opts)
	err := br.AddBlock(*genesisBlock)
	assert.NoError(t, err)
}

func TestBlockRepository_GetBlockByHeight(t *testing.T) {
	dbPath := "./.db"

	genesisConfFilePath := build.Default.GOPATH + "/src/github.com/it-chain/it-chain-Engine/.it-chain/genesisconf/GenesisBlockConfig.json"

	validator := new(impl.DefaultValidator)
	db := leveldbwrapper.CreateNewDB(dbPath)
	defer func() {
		db.Close()
		os.RemoveAll(dbPath)
	}()
	opts := map[string]interface{}{
		"db_path": dbPath,
	}
	br := NewBlockRepository(db, *validator, opts)
	genesisBlock, _ := blockchain.CreateGenesisBlock(genesisConfFilePath)
	br.AddBlock(*genesisBlock)
	retrievedBlock, err := br.GetBlockByHeight(0)
	assert.NoError(t, err)
	assert.Equal(t, genesisBlock, retrievedBlock)
}

func TestBlockRepository_GetBlockBySeal(t *testing.T) {
	dbPath := "./.db"

	genesisConfFilePath := build.Default.GOPATH + "/src/github.com/it-chain/it-chain-Engine/.it-chain/genesisconf/GenesisBlockConfig.json"

	validator := new(impl.DefaultValidator)
	db := leveldbwrapper.CreateNewDB(dbPath)
	defer func() {
		db.Close()
		os.RemoveAll(dbPath)
	}()
	opts := map[string]interface{}{
		"db_path": dbPath,
	}
	br := NewBlockRepository(db, *validator, opts)
	genesisBlock, _ := blockchain.CreateGenesisBlock(genesisConfFilePath)
	br.AddBlock(*genesisBlock)
	retrievedBlock, err := br.GetBlockBySeal(genesisBlock.Seal)
	assert.NoError(t, err)
	assert.Equal(t, genesisBlock, retrievedBlock)
}

func TestBlockRepository_GetLastBlock(t *testing.T) {
	dbPath := "./.db"

	genesisConfFilePath := build.Default.GOPATH + "/src/github.com/it-chain/it-chain-Engine/.it-chain/genesisconf/GenesisBlockConfig.json"

	validator := new(impl.DefaultValidator)
	db := leveldbwrapper.CreateNewDB(dbPath)
	defer func() {
		db.Close()
		os.RemoveAll(dbPath)
	}()
	opts := map[string]interface{}{
		"db_path": dbPath,
	}
	br := NewBlockRepository(db, *validator, opts)
	genesisBlock, _ := blockchain.CreateGenesisBlock(genesisConfFilePath)
	br.AddBlock(*genesisBlock)
	retrievedBlock, err := br.GetLastBlock()
	assert.NoError(t, err)
	assert.Equal(t, genesisBlock, retrievedBlock)
}

//Todo: GetTransactionByTxID, GetBlockByTxID 테스트 코드 작성(주석작성자:junk-sound)
