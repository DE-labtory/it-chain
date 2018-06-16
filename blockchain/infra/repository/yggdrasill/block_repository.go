package leveldb

import (
	"github.com/it-chain/leveldb-wrapper"
	"github.com/it-chain/yggdrasill"
	"github.com/it-chain/yggdrasill/common"
)

type BlockRepository = yggdrasill.Yggdrasill

func NewBlockRepository(dbpath string, validator common.Validator, opts map[string]interface{}) (*BlockRepository, error) {
	db := leveldbwrapper.CreateNewDB(dbpath)

	return yggdrasill.NewYggdrasill(db, validator, opts)
}
