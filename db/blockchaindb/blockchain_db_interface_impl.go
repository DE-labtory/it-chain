package blockchaindb

import (
	"github.com/spf13/viper"
	"it-chain/db/blockchaindb/blockchainleveldb"
	"it-chain/service/blockchain"
	"fmt"
)

type blockchainDBImpl struct {
	dbType string
	db     BlockChainDB
}

func CreateNewBlockchainDB(dbPath string) *blockchainDBImpl {
	dbType := viper.GetString("database.type")
	var db BlockChainDB
	switch dbType {
	case "leveldb":
		db = blockchainleveldb.CreateNewBlockchainLevelDB(dbPath)
		break
	default :
		panic(fmt.Sprint("Unsupported db type"))
	}
	return &blockchainDBImpl{dbType: dbType, db: db}
}

func (b *blockchainDBImpl) Close() {
	b.db.Close()
}

func (b *blockchainDBImpl) AddBlock(block *blockchain.Block) error {
	return b.db.AddBlock(block)
}

func (b *blockchainDBImpl) AddUnconfirmedBlock(block *blockchain.Block) error {
	return b.db.AddUnconfirmedBlock(block)
}

func (b *blockchainDBImpl) GetBlockByNumber(blockNumber uint64) (*blockchain.Block, error) {
	return b.db.GetBlockByNumber(blockNumber)
}

func (b *blockchainDBImpl) GetBlockByHash(hash string) (*blockchain.Block, error) {
	return b.db.GetBlockByHash(hash)
}

func (b *blockchainDBImpl) GetLastBlock() (*blockchain.Block, error) {
	return b.db.GetLastBlock()
}

func (b *blockchainDBImpl) GetTransactionByTxID(txid string) (*blockchain.Transaction, error) {
	return b.db.GetTransactionByTxID(txid)
}

func (b *blockchainDBImpl) GetBlockByTxID(txid string) (*blockchain.Block, error) {
	return b.db.GetBlockByTxID(txid)
}