package api

import (
	"testing"

	"fmt"

	"github.com/it-chain/it-chain-Engine/blockchain/domain/model/block"
	"github.com/it-chain/it-chain-Engine/blockchain/domain/repository"
	"github.com/it-chain/yggdrasill"
	ygg "github.com/it-chain/yggdrasill"
	b "github.com/it-chain/yggdrasill/block"
	tx "github.com/it-chain/yggdrasill/transaction"
	"github.com/magiconair/properties/assert"
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
	Yggdrasill *blockchainleveldb.YggDrasill
}

func (br *BlockRepositoryImpl) Close() {

}
func (br *BlockRepositoryImpl) AddBlock(block block.Block) error {
	return nil
}
func (br *BlockRepositoryImpl) GetLastBlock(block block.Block) {

}
func (br *BlockRepositoryImpl) GetTransactionsById(id int) { // block과 관련된 정보 조회 예시

}
func NewBlockRepository(yggdrasill *blockchainleveldb.YggDrasill) repository.BlockRepository {
	return &BlockRepositoryImpl{
		Yggdrasill: yggdrasill,
	}
}

func TestBlockApi_AddBlock_DefaultBlock_basic(t *testing.T) {
	// given
	dbPath := "./db"

	y := ygg.NewYggdrasil(dbPath, nil)
	br := NewBlockRepository(y)
	bApi := NewBlockApi(br)
	block1 := BlockImpl{
		Header: b.BlockHeader{
			Height:    0,
			CreatorID: "zf",
		},
	}
	block2 := BlockImpl{}

	// when
	err1 := bApi.AddBlock(&block1)
	err2 := bApi.AddBlock(&block2)

	// then
	if err1 != nil {
		fmt.Print(err1.Error())
	}
	if err2 != nil {
		fmt.Print(err2.Error())
	}
	assert.Equal(t, block1.GetHeight(), block2.GetHeight())
	assert.Equal(t, uint64(0), block1.GetHeight())
	assert.Equal(t, uint64(0), block2.GetHeight())
	assert.Equal(t, "zf", block1.Header.CreatorID)
}
