package yggdrasill

import (
	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/yggdrasill"
	"github.com/it-chain/yggdrasill/common"
	"github.com/it-chain/yggdrasill/impl"
)

type BlockRepository struct {
	*yggdrasill.BlockStorage
	Creator string
}

type Block = impl.DefaultBlock

func NewBlockRepository(dbpath string, opts map[string]interface{}, creator string) (*BlockRepository, error) {
	// Use default validator
	var validator common.Validator
	validator = new(impl.DefaultValidator)

	db := leveldbwrapper.CreateNewDB(dbpath)

	storage, err := yggdrasill.NewBlockStorage(db, validator, opts)

	if err != nil {
		return nil, err
	}

	return &BlockRepository{storage, creator}, nil
}

func (b *BlockRepository) NewEmptyBlock() (*impl.DefaultBlock, error) {
	lastBlock := &Block{}
	err := b.GetLastBlock(lastBlock)
	if err != nil {
		return nil, err
	}

	prevSeal := lastBlock.GetSeal()
	height := lastBlock.GetHeight() + 1 // TODO: correct?
	creator := []byte(b.Creator)

	return impl.NewEmptyBlock(prevSeal, height, creator), nil
}

func (b *BlockRepository) GetBlockCreator() string {
	return b.Creator
}
