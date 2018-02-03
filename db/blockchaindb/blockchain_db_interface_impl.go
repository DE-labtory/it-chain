package blockchaindb

import (
	"github.com/spf13/viper"
	"it-chain/db/blockchaindb/blockchainleveldb"
	"it-chain/service/blockchain"
	"fmt"
)

type BlockchainDBImpl struct {
	dbType string
	db     BlockChainDB
}

func CreateNewBlockchainDB(dbPath string) *BlockchainDBImpl {
	dbType := viper.GetString("database.type")
	var db BlockChainDB
	switch dbType {
	case "leveldb":
		db = blockchainleveldb.CreateNewBlockchainLevelDB(dbPath)
		break
	default :
		panic(fmt.Sprint("Unsupported db type"))
	}
	return &BlockchainDBImpl{dbType: dbType, db: db}
}

func (b *BlockchainDBImpl) Close() {
	b.db.Close()
}

func (b *BlockchainDBImpl) AddBlock(block *blockchain.Block) error {
	return b.db.AddBlock(block)
}

func (b *BlockchainDBImpl) AddUnconfirmedBlock(block *blockchain.Block) error {
	return b.db.AddUnconfirmedBlock(block)
}

func (b *BlockchainDBImpl) GetBlockByNumber(blockNumber uint64) (*blockchain.Block, error) {
	return b.db.GetBlockByNumber(blockNumber)
}

func (b *BlockchainDBImpl) GetBlockByHash(hash string) (*blockchain.Block, error) {
	return b.db.GetBlockByHash(hash)
}

func (b *BlockchainDBImpl) GetLastBlock() (*blockchain.Block, error) {
	return b.db.GetLastBlock()
}

func (b *BlockchainDBImpl) GetTransactionByTxID(txid string) (*blockchain.Transaction, error) {
	return b.db.GetTransactionByTxID(txid)
}

func (b *BlockchainDBImpl) GetBlockByTxID(txid string) (*blockchain.Block, error) {
	return b.db.GetBlockByTxID(txid)
}