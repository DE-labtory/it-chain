package repository

import (
	"log"

	//"github.com/it-chain/it-chain-Engine/blockchain/domain/model/block"
	"github.com/it-chain/it-chain-Engine/blockchain/domain/model/transaction"
	"github.com/it-chain/leveldb-wrapper/key_value_db"
	"github.com/it-chain/yggdrasill"
	"github.com/it-chain/yggdrasill/common"
	"github.com/it-chain/yggdrasill/impl"
)

type BlockRepository struct {
	yggdrasill *yggdrasill.Yggdrasill
}


func NewBlockRepository(keyValueDB key_value_db.KeyValueDB, validator impl.DefaultValidator, opts map[string]interface{}) *BlockRepository {
	ygg, err := yggdrasill.NewYggdrasill(keyValueDB, &validator, opts)

	if err != nil {
		log.Fatal(err.Error())
	}

	return &BlockRepository{
		yggdrasill: ygg,
	}
}

func (br BlockRepository) Close() {
	br.yggdrasill.Close()
}

func (br BlockRepository) AddBlock(block common.Block) error {
	err := br.yggdrasill.AddBlock(block)
	if err != nil {
		return err
	}
	return nil
}

// Issue : func GetBlock~
// 현재 repository가 it-chain-engine에 속하기 때문에 ReturnBlock의 type을 it-chain-Engine/blockchain/domain/model/block.Block으로 하고 싶지만,
// 그렇게 되면 common.Block 인터페이스에 속하지 않는 것 같음. block.Block에 필요한 메소드를 다시 구현하기보다는 현재는 impl.Defaultblock 사용중(주석작성자 GitID:junk-sound).

func (br BlockRepository) GetBlockByNumber(blockNumber uint64) (*impl.DefaultBlock, error) {
	var ReturnBlock impl.DefaultBlock
	err := br.yggdrasill.GetBlockByHeight(&ReturnBlock, blockNumber)
	if err != nil {
		return nil, err
	}
	return &ReturnBlock, err
}


func (br BlockRepository) GetBlockBySeal(seal []byte) (*impl.DefaultBlock, error) {
	var ReturnBlock impl.DefaultBlock
	err := br.yggdrasill.GetBlockBySeal(&ReturnBlock, seal)
	if err != nil {
		return nil, err
	}
	return &ReturnBlock, err
}

func (br BlockRepository) GetLastBlock() (*impl.DefaultBlock, error) {
	var ReturnBlock impl.DefaultBlock
	err := br.yggdrasill.GetLastBlock(&ReturnBlock)
	if err != nil {
		return nil,err
	}
	return &ReturnBlock, err
}

func (br BlockRepository) GetTransactionByTxID(txid string) (*transaction.Trasaction, error) {
	var ReturnTx impl.DefaultTransaction
	err := br.yggdrasill.GetTransactionByTxID(&ReturnTx, txid)
	if err != nil {
		return nil, err
	}
	return &ReturnTx, err
}

func (br BlockRepository) GetBlockByTxID(txid string) (*impl.DefaultBlock, error) {
	var ReturnBlock impl.DefaultBlock
	err := br.yggdrasill.GetBlockByTxID(&ReturnBlock, txid)
	if err != nil {
		return nil, err
	}
	return &ReturnBlock, err
}
