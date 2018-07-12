package yggdrasill

//import (
//	"github.com/it-chain/it-chain-Engine/blockchain"
//	"github.com/it-chain/leveldb-wrapper"
//	"github.com/it-chain/yggdrasill"
//)
//
//type BlockRepository struct {
//	*yggdrasill.BlockStorage
//	Creator string
//}
//
//func NewBlockRepository(dbpath string, opts map[string]interface{}, creator string) (*BlockRepository, error) {
//	// Use default validator
//	var validator blockchain.Validator
//	validator = new(blockchain.DefaultValidator)
//
//	db := leveldbwrapper.CreateNewDB(dbpath)
//
//	storage, err := yggdrasill.NewBlockStorage(db, validator, opts)
//
//	if err != nil {
//		return nil, err
//	}
//
//	return &BlockRepository{storage, creator}, nil
//}
//
//func (b *BlockRepository) NewEmptyBlock() (blockchain.Block, error) {
//	lastBlock := &blockchain.DefaultBlock{}
//	err := b.GetLastBlock(lastBlock)
//	if err != nil {
//		return nil, err
//	}
//
//	prevSeal := lastBlock.GetSeal()
//	height := lastBlock.GetHeight() + 1 // TODO: correct?
//	creator := []byte(b.Creator)
//
//	return blockchain.NewEmptyBlock(prevSeal, height, creator), nil
//}
//
//func (b *BlockRepository) GetBlockCreator() string {
//	return b.Creator
//}
//
//}
