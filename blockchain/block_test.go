package blockchain

import (
	"os"
	"testing"

	"time"

	"encoding/json"
	"io/ioutil"

	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/yggdrasill/impl"
	"github.com/stretchr/testify/assert"
)

func TestBlockRepository_AddBlock(t *testing.T) {
	GenesisConfFilePath := "./GenesisBlockConfig.json"
	defer os.Remove(GenesisConfFilePath)
	GenesisBlockConfigJson := []byte(`{
								  "Seal":[],
								  "PrevSeal":[],
								  "Height":0,
								  "TxList":[],
								  "TxSeal":[],
								  "TimeStamp":"0001-01-01T00:00:00-00:00",
								  "Creator":[]
								}`)
	var tempJson impl.DefaultBlock
	_ = json.Unmarshal(GenesisBlockConfigJson, &tempJson)
	GenesisBlockConfigByte, _ := json.Marshal(tempJson)
	_ = ioutil.WriteFile(GenesisConfFilePath, GenesisBlockConfigByte, 0644)

	dbPath := "./.db"
	db := leveldbwrapper.CreateNewDB(dbPath)
	defer func() {
		db.Close()
		os.RemoveAll(dbPath)
	}()

	genesisBlock, _ := CreateGenesisBlock(GenesisConfFilePath)
	validator := new(impl.DefaultValidator)
	opts := map[string]interface{}{
		"db_path": dbPath,
	}
	br := NewBlockRepository(db, *validator, opts)
	err := br.AddBlock(*genesisBlock)
	assert.NoError(t, err)

}

func TestBlockRepository_GetBlockByHeight(t *testing.T) {
	GenesisConfFilePath := "./GenesisBlockConfig.json"
	defer os.Remove(GenesisConfFilePath)
	GenesisBlockConfigJson := []byte(`{
								  "Seal":[],
								  "PrevSeal":[],
								  "Height":0,
								  "TxList":[],
								  "TxSeal":[],
								  "TimeStamp":"0001-01-01T00:00:00-00:00",
								  "Creator":[]
								}`)
	var tempJson impl.DefaultBlock
	_ = json.Unmarshal(GenesisBlockConfigJson, &tempJson)
	GenesisBlockConfigByte, _ := json.Marshal(tempJson)
	_ = ioutil.WriteFile(GenesisConfFilePath, GenesisBlockConfigByte, 0644)

	dbPath := "./.db"
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
	genesisBlock, _ := CreateGenesisBlock(GenesisConfFilePath)
	br.AddBlock(*genesisBlock)
	retrievedBlock, err := br.GetBlockByHeight(0)
	assert.NoError(t, err)
	assert.Equal(t, genesisBlock, retrievedBlock)
}

func TestBlockRepository_GetBlockBySeal(t *testing.T) {
	GenesisConfFilePath := "./GenesisBlockConfig.json"
	defer os.Remove(GenesisConfFilePath)
	GenesisBlockConfigJson := []byte(`{
								  "Seal":[],
								  "PrevSeal":[],
								  "Height":0,
								  "TxList":[],
								  "TxSeal":[],
								  "TimeStamp":"0001-01-01T00:00:00-00:00",
								  "Creator":[]
								}`)
	var tempJson impl.DefaultBlock
	_ = json.Unmarshal(GenesisBlockConfigJson, &tempJson)
	GenesisBlockConfigByte, _ := json.Marshal(tempJson)
	_ = ioutil.WriteFile(GenesisConfFilePath, GenesisBlockConfigByte, 0644)

	dbPath := "./.db"
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
	genesisBlock, _ := CreateGenesisBlock(GenesisConfFilePath)
	br.AddBlock(*genesisBlock)
	retrievedBlock, err := br.GetBlockBySeal(genesisBlock.Seal)
	assert.NoError(t, err)
	assert.Equal(t, genesisBlock, retrievedBlock)
}

func TestBlockRepository_GetLastBlock(t *testing.T) {
	GenesisConfFilePath := "./GenesisBlockConfig.json"
	defer os.Remove(GenesisConfFilePath)
	GenesisBlockConfigJson := []byte(`{
								  "Seal":[],
								  "PrevSeal":[],
								  "Height":0,
								  "TxList":[],
								  "TxSeal":[],
								  "TimeStamp":"0001-01-01T00:00:00-00:00",
								  "Creator":[]
								}`)
	var tempJson impl.DefaultBlock
	_ = json.Unmarshal(GenesisBlockConfigJson, &tempJson)
	GenesisBlockConfigByte, _ := json.Marshal(tempJson)
	_ = ioutil.WriteFile(GenesisConfFilePath, GenesisBlockConfigByte, 0644)

	dbPath := "./.db"
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
	genesisBlock, _ := CreateGenesisBlock(GenesisConfFilePath)
	br.AddBlock(*genesisBlock)
	retrievedBlock, err := br.GetLastBlock()
	assert.NoError(t, err)
	assert.Equal(t, genesisBlock, retrievedBlock)
}

func TestBlockRepositoryImpl_GetBlockByTxID(t *testing.T) {
	GenesisConfFilePath := "./GenesisBlockConfig.json"
	defer os.Remove(GenesisConfFilePath)
	GenesisBlockConfigJson := []byte(`{
								  "Seal":[],
								  "PrevSeal":[],
								  "Height":0,
								  "TxList":[],
								  "TxSeal":[],
								  "TimeStamp":"0001-01-01T00:00:00-00:00",
								  "Creator":[]
								}`)
	var tempJson impl.DefaultBlock
	_ = json.Unmarshal(GenesisBlockConfigJson, &tempJson)
	GenesisBlockConfigByte, _ := json.Marshal(tempJson)
	_ = ioutil.WriteFile(GenesisConfFilePath, GenesisBlockConfigByte, 0644)

	dbPath := "./.db"
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
	genesisBlock, _ := CreateGenesisBlock(GenesisConfFilePath)
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
	GenesisConfFilePath := "./GenesisBlockConfig.json"
	defer os.Remove(GenesisConfFilePath)
	GenesisBlockConfigJson := []byte(`{
								  "Seal":[],
								  "PrevSeal":[],
								  "Height":0,
								  "TxList":[],
								  "TxSeal":[],
								  "TimeStamp":"0001-01-01T00:00:00-00:00",
								  "Creator":[]
								}`)
	var tempJson impl.DefaultBlock
	_ = json.Unmarshal(GenesisBlockConfigJson, &tempJson)
	GenesisBlockConfigByte, _ := json.Marshal(tempJson)
	_ = ioutil.WriteFile(GenesisConfFilePath, GenesisBlockConfigByte, 0644)

	dbPath := "./.db"
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
	genesisBlock, _ := CreateGenesisBlock(GenesisConfFilePath)
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
