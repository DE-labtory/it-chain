package yggdrasill

import (
	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/yggdrasill"
	"github.com/it-chain/yggdrasill/common"
	"github.com/it-chain/yggdrasill/impl"
)

type BlockRepository = yggdrasill.BlockStorage

type Block = impl.DefaultBlock

func NewBlockRepository(dbpath string, opts map[string]interface{}) (*BlockRepository, error) {
	// Use default validator
	var validator common.Validator
	validator = new(impl.DefaultValidator)

	db := leveldbwrapper.CreateNewDB(dbpath)

	return yggdrasill.NewBlockStorage(db, validator, opts)
}
