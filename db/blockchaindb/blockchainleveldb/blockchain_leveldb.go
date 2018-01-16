package blockchainleveldb

import (
	"it-chain/db/leveldbhelper"
	"it-chain/service/blockchain"
	"it-chain/common"
	"fmt"
)

const (
	BLOCK_HASH_DB = "block_hash"
	BLOCK_NUMBER_DB = "block_number"
	TRANSACTION_DB = "transaction"
)

type BlockchainLevelDB struct {
	DBProvider *leveldbhelper.DBProvider
}

func CreateNewBlockchainLevelDB(levelDBPath string) *BlockchainLevelDB {
	levelDBProvider := leveldbhelper.CreateNewDBProvider(levelDBPath)
	return &BlockchainLevelDB{levelDBProvider}
}

func (l *BlockchainLevelDB) Close() {
	l.DBProvider.Close()
}

func (l *BlockchainLevelDB) AddBlock(block *blockchain.Block) error {
	blockHashDB := l.DBProvider.GetDBHandle(BLOCK_HASH_DB)
	blockNumberDB := l.DBProvider.GetDBHandle(BLOCK_NUMBER_DB)

	serializedBlock, err := common.Serialize(block)
	if err != nil {
		return err
	}

	err = blockNumberDB.Put([]byte(fmt.Sprint(block.Header.Number)), serializedBlock, true)
	if err != nil {
		return err
	}

	err = blockHashDB.Put([]byte(block.Header.BlockHash), serializedBlock, true)
	if err != nil {
		return err
	}

	return nil
}


func (l *BlockchainLevelDB) GetBlockByNumber(blockNumber uint64) (*blockchain.Block, error) {
	blockNumberDB := l.DBProvider.GetDBHandle(BLOCK_NUMBER_DB)

	serializedBlock, err := blockNumberDB.Get([]byte(fmt.Sprint(blockNumber)))
	if err != nil {
		return nil, err
	}

	block := &blockchain.Block{}
	err = common.Deserialize(serializedBlock, block)
	if err != nil {
		return nil, err
	}

	return block, err
}

func (l *BlockchainLevelDB) GetBlockByHash(hash string) (*blockchain.Block, error) {
	blockHashDB := l.DBProvider.GetDBHandle(BLOCK_HASH_DB)

	serializedBlock, err := blockHashDB.Get([]byte(hash))
	if err != nil {
		return nil, err
	}

	block := &blockchain.Block{}
	err = common.Deserialize(serializedBlock, block)
	if err != nil {
		return nil, err
	}

	return block, err
}