package blockchain

import (
	"go/build"
	"os"
	"testing"

	"time"

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

	genesisBlock, _ := CreateGenesisBlock(genesisConfFilePath)
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
	genesisBlock, _ := CreateGenesisBlock(genesisConfFilePath)
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
	genesisBlock, _ := CreateGenesisBlock(genesisConfFilePath)
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
	genesisBlock, _ := CreateGenesisBlock(genesisConfFilePath)
	br.AddBlock(*genesisBlock)
	retrievedBlock, err := br.GetLastBlock()
	assert.NoError(t, err)
	assert.Equal(t, genesisBlock, retrievedBlock)
}

func TestBlockRepositoryImpl_GetBlockByTxID(t *testing.T) {
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
	genesisBlock, _ := CreateGenesisBlock(genesisConfFilePath)
	tx := &impl.DefaultTransaction{
		PeerID:    "p01",
		ID:        "tx01",
		Status:    0,
		Timestamp: time.Now().Round(0),
		TxData: &impl.TxData{
			Jsonrpc: "jsonRPC01",
			Method:  "invoke",
			Params: impl.Params{
				Type:     0,
				Function: "function01",
				Args:     []string{"arg1", "arg2"},
			},
			ID: "txdata01",
		},
		Signature: nil,
	}
	genesisBlock.PutTx(tx)
	err := br.AddBlock(*genesisBlock)
	assert.NoError(t, err)
	retrievedBlock, err := br.GetBlockByTxID("tx01")
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), retrievedBlock.Height)

}

func TestBlockRepositoryImpl_GetTransactionByTxID(t *testing.T) {
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
	genesisBlock, _ := CreateGenesisBlock(genesisConfFilePath)
	tx := &impl.DefaultTransaction{
		PeerID:    "p01",
		ID:        "tx01",
		Status:    0,
		Timestamp: time.Now().Round(0),
		TxData: &impl.TxData{
			Jsonrpc: "jsonRPC01",
			Method:  "invoke",
			Params: impl.Params{
				Type:     0,
				Function: "function01",
				Args:     []string{"arg1", "arg2"},
			},
			ID: "txdata01",
		},
		Signature: nil,
	}
	genesisBlock.PutTx(tx)
	err := br.AddBlock(*genesisBlock)
	assert.NoError(t, err)
	retrievedTx, err := br.GetTransactionByTxID("tx01")
	assert.Equal(t, "tx01", retrievedTx.ID)
}
