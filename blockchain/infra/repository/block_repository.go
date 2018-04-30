package repository

import (
	"github.com/it-chain/leveldb-wrapper/key_value_db"
	"github.com/it-chain/yggdrasill"
	"github.com/it-chain/yggdrasill/validator"
	"github.com/it-chain/it-chain-Engine/blockchain/domain/model/block"
	"github.com/it-chain/it-chain-Engine/blockchain/domain/model/transaction"
)

type BlockRepository struct {
	yggdrasill *blockchaindb.Yggdrasill
}

func NewBlockRepository(keyValueDB key_value_db.KeyValueDB, validator validator.Validator, opts map[string]interface{}) *BlockRepository {
	ygg := blockchaindb.NewYggdrasill(keyValueDB, validator, opts);
	return &BlockRepository{
		yggdrasill: ygg,
	}
}

// TODO:
// Yggdrasill wrapper 함수 만들기
func (br BlockRepository) Close() {}
func (br BlockRepository) AddBlock(block block.Block) error {}
func (br BlockRepository) GetBlockByNumber(blockNumber uint64) (block.Block, error) {}
func (br BlockRepository) GetBlockByHash(hash string) (block.Block, error) {}
func (br BlockRepository) GetLastBlock() (block.Block, error) {}
func (br BlockRepository) GetTransactionByTxID(txid string) (transaction.Trasaction, error) {}
func (br BlockRepository) GetBlockByTxID(txid string) (block.Block, error) {}