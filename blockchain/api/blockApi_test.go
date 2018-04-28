package api

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/blockchain/domain/model/block"
	"github.com/it-chain/it-chain-Engine/blockchain/domain/repository"
	"github.com/it-chain/yggdrasill"
	b "github.com/it-chain/yggdrasill/block"
	tx "github.com/it-chain/yggdrasill/transaction"
)

/*
	DefaultBlock type BlockImpl
*/
type BlockImpl b.DefaultBlock

func (b *BlockImpl) PutTransaction(transaction tx.Transaction) {

}

func (b *BlockImpl) FindTransactionIndexByHash(txHash string) {

}

func (b *BlockImpl) Serialize() ([]byte, error) {
	return nil, nil
}

func (b *BlockImpl) GenerateHash() error {
	return nil
}

func (b *BlockImpl) GetHash() string {
	return ""
}

func (b *BlockImpl) GetTransactions() []tx.Transaction {
	return nil
}

func (b *BlockImpl) GetHeight() uint64 {
	return b.Header.Height
}

func (b *BlockImpl) IsPrev(serializedBlock []byte) bool {
	return false
}

/*
	BlockRepositoryImpl
*/
type BlockRepositoryImpl struct {
	Yggdrasill *blockchaindb.Yggdrasill
}

func (br *BlockRepositoryImpl) Close() {

}
func (br *BlockRepositoryImpl) AddBlock(block block.Block) error {
	return nil
}
func (br *BlockRepositoryImpl) GetLastBlock() block.Block {
	return nil
}
func (br *BlockRepositoryImpl) GetTransactionsById(id string) tx.Transaction {
	return nil
}
func NewBlockRepository(yggdrasill *blockchaindb.Yggdrasill) repository.BlockRepository {
	return &BlockRepositoryImpl{
		Yggdrasill: yggdrasill,
	}
}

func TestBlockApi_AddBlock_DefaultBlock_basic(t *testing.T) {

}
