package blockchain

import (
	"github.com/it-chain/yggdrasill/impl"
	"github.com/it-chain/yggdrasill"
	"github.com/it-chain/leveldb-wrapper/key_value_db"
	"log"
)

type Block = impl.DefaultBlock


type BlockRepository interface {
	Close()
	AddBlock(block Block) error
	GetBlockByHeight(blockHeight uint64) (*Block, error)
	GetBlockBySeal(seal []byte) (*Block, error)
	GetLastBlock() (*Block, error)
	GetTransactionByTxID(txid string) (*Transaction, error)
	GetBlockByTxID(txid string) (*Block, error)
}

type BlockRepositoryImpl struct {
	yggdrasill *yggdrasill.Yggdrasill
}

func NewBlockRepository(keyValueDB key_value_db.KeyValueDB, validator impl.DefaultValidator, opts map[string]interface{}) BlockRepository {
	ygg, err := yggdrasill.NewYggdrasill(keyValueDB, &validator, opts)

	if err != nil {
		log.Fatal(err.Error())
	}
	return &BlockRepositoryImpl{
		yggdrasill: ygg,
	}
}

func (br *BlockRepositoryImpl) Close() {
	br.yggdrasill.Close()
}

func (br *BlockRepositoryImpl) AddBlock(block Block) error {
	var yggBlock = impl.DefaultBlock(block)
	err := br.yggdrasill.AddBlock(&yggBlock)
	if err != nil {
		return err
	}
	return nil
}

func (br *BlockRepositoryImpl) GetBlockByHeight(blockHeight uint64) (*Block, error) {
	var retrievedBlock impl.DefaultBlock
	err := br.yggdrasill.GetBlockByHeight(&retrievedBlock, blockHeight)
	if err != nil {
		return nil, err
	}
	return &retrievedBlock, err
}

func (br *BlockRepositoryImpl) GetBlockBySeal(seal []byte) (*Block, error) {
	var retrievedBlock impl.DefaultBlock
	err := br.yggdrasill.GetBlockBySeal(&retrievedBlock, seal)
	if err != nil {
		return nil, err
	}
	return &retrievedBlock, err
}

func (br *BlockRepositoryImpl) GetLastBlock() (*Block, error) {
	var retrievedBlock impl.DefaultBlock
	err := br.yggdrasill.GetLastBlock(&retrievedBlock)
	if err != nil {
		return nil, err
	}
	return &retrievedBlock, err
}

func (br *BlockRepositoryImpl) GetTransactionByTxID(txid string) (*Transaction, error) {
	var retrievedTx impl.DefaultTransaction
	err := br.yggdrasill.GetTransactionByTxID(&retrievedTx, txid)
	if err != nil {
		return nil, err
	}
	return &retrievedTx, err
}

func (br *BlockRepositoryImpl) GetBlockByTxID(txid string) (*Block, error) {
	var retrievedBlock impl.DefaultBlock
	err := br.yggdrasill.GetBlockByTxID(&retrievedBlock, txid)
	if err != nil {
		return nil, err
	}
	return &retrievedBlock, err
}
