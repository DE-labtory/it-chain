package repository

import (
	"log"

	"github.com/it-chain/it-chain-Engine/blockchain/domain/model/block"
	"github.com/it-chain/it-chain-Engine/blockchain/domain/model/transaction"
	"github.com/it-chain/leveldb-wrapper/key_value_db"
	"github.com/it-chain/yggdrasill"
	"github.com/it-chain/yggdrasill/impl"
)

type BlockRepository struct {
	yggdrasill *yggdrasill.Yggdrasill
}

func NewBlockRepository(keyValueDB key_value_db.KeyValueDB, validator impl.DefaultValidator, opts map[string]interface{}) *BlockRepository {
	ygg, err := yggdrasill.NewYggdrasill(keyValueDB, validator, opts)

	if err != nil {
		log.Fatal(err.Error())
	}

	return &BlockRepository{
		yggdrasill: ygg,
	}
}

// TODO:
// Yggdrasill wrapper 함수 만들기
func (br BlockRepository) Close() {

}

func (br BlockRepository) AddBlock(block block.Block) error {

}

func (br BlockRepository) GetBlockByNumber(blockNumber uint64) (block.Block, error) {

}

func (br BlockRepository) GetBlockByHash(hash string) (block.Block, error) {

}

func (br BlockRepository) GetLastBlock() (block.Block, error) {

}

func (br BlockRepository) GetTransactionByTxID(txid string) (transaction.Trasaction, error) {

}

func (br BlockRepository) GetBlockByTxID(txid string) (block.Block, error) {

}
