package blockchainleveldb

import (
	"github.com/spf13/viper"
	"fmt"
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/blockchain/domain/model"
	"github.com/it-chain/it-chain-Engine/common/leveldbhelper"
)

var logger = common.GetLogger("blockchain_leveldb.go")

const (
	BLOCK_HASH_DB = "block_hash"
	BLOCK_NUMBER_DB = "block_number"
	UNCONFIRMED_BLOCK_DB = "unconfirmed_block"
	TRANSACTION_DB = "transaction"
	UTIL_DB = "util"

	LAST_BLOCK_KEY = "last_block"
)

type BlockchainLevelDB struct {
	DBProvider *leveldbhelper.DBProvider
}

func CreateNewBlockchainLevelDB(levelDBPath string) *BlockchainLevelDB {
	if levelDBPath == "" {
		levelDBPath = viper.GetString("database.leveldb.defaultPath")
	}
	levelDBProvider := leveldbhelper.CreateNewDBProvider(levelDBPath)

	return &BlockchainLevelDB{levelDBProvider}
}

func (l *BlockchainLevelDB) Close() {
	l.DBProvider.Close()
}

func (l *BlockchainLevelDB) AddBlock(block *model.Block) error {
	blockHashDB := l.DBProvider.GetDBHandle(BLOCK_HASH_DB)
	blockNumberDB := l.DBProvider.GetDBHandle(BLOCK_NUMBER_DB)
	transactionDB := l.DBProvider.GetDBHandle(TRANSACTION_DB)
	unconfirmedDB := l.DBProvider.GetDBHandle(UNCONFIRMED_BLOCK_DB)
	utilDB := l.DBProvider.GetDBHandle(UTIL_DB)

	serializedBlock, err := common.Serialize(block)
	if err != nil {
		return err
	}

	err = blockHashDB.Put([]byte(block.BlockHeader.BlockHash), serializedBlock, true)
	if err != nil {
		return err
	}

	err = blockNumberDB.Put([]byte(fmt.Sprint(block.BlockHeader.Number)), []byte(block.BlockHeader.BlockHash), true)
	if err != nil {
		return err
	}

	err = utilDB.Put([]byte(LAST_BLOCK_KEY), serializedBlock, true)
	if err != nil {
		return err
	}

	for _, tx := range block.BlockData.Transactions {
		serializedTx, err := common.Serialize(tx)
		if err != nil {
			return err
		}

		err = transactionDB.Put([]byte(tx.TransactionID), serializedTx, true)
		if err != nil {
			return err
		}

		err = utilDB.Put([]byte(tx.TransactionID), []byte(block.BlockHeader.BlockHash), true)
		if err != nil {
			return err
		}
	}

	err = unconfirmedDB.Delete([]byte(block.BlockHeader.BlockHash), true)
	if err != nil {
		return err
	}

	return nil
}

func (l *BlockchainLevelDB) AddUnconfirmedBlock(block *model.Block) error {
	unconfirmedDB := l.DBProvider.GetDBHandle(UNCONFIRMED_BLOCK_DB)

	serializedBlock, err := common.Serialize(block)
	if err != nil {
		return err
	}

	err = unconfirmedDB.Put([]byte(block.BlockHeader.BlockHash), serializedBlock, true)
	if err != nil {
		return err
	}

	return nil
}

func (l *BlockchainLevelDB) GetBlockByNumber(blockNumber uint64) (*model.Block, error) {
	blockNumberDB := l.DBProvider.GetDBHandle(BLOCK_NUMBER_DB)

	blockHash, err := blockNumberDB.Get([]byte(fmt.Sprint(blockNumber)))
	if err != nil {
		return nil, err
	}

	return l.GetBlockByHash(string(blockHash))
}

func (l *BlockchainLevelDB) GetBlockByHash(hash string) (*model.Block, error) {
	blockHashDB := l.DBProvider.GetDBHandle(BLOCK_HASH_DB)

	serializedBlock, err := blockHashDB.Get([]byte(hash))
	if err != nil {
		return nil, err
	}

	block := &model.Block{}
	err = common.Deserialize(serializedBlock, block)
	if err != nil {
		return nil, err
	}

	return block, err
}

func (l *BlockchainLevelDB) GetLastBlock() (*model.Block, error) {
	utilDB := l.DBProvider.GetDBHandle(UTIL_DB)

	serializedBlock, err := utilDB.Get([]byte(LAST_BLOCK_KEY))

	if err != nil{
		return nil, err
	}

	if serializedBlock == nil{
		return nil, nil
	}

	block := &model.Block{}
	err = common.Deserialize(serializedBlock, block)
	if err != nil {
		return nil, err
	}

	return block, err
}

func (l *BlockchainLevelDB) GetTransactionByTxID(txid string) (*model.Transaction, error) {
	transactionDB := l.DBProvider.GetDBHandle(TRANSACTION_DB)

	serializedTX, err := transactionDB.Get([]byte(txid))
	if err != nil {
		return nil, err
	}

	transaction := &model.Transaction{}
	err = common.Deserialize(serializedTX, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, err
}

func (l *BlockchainLevelDB) GetBlockByTxID(txid string) (*model.Block, error) {
	utilDB := l.DBProvider.GetDBHandle(UTIL_DB)

	blockHash, err := utilDB.Get([]byte(txid))
	if err != nil {
		return nil, err
	}

	return l.GetBlockByHash(string(blockHash))
}