package leveldb

import (
	"github.com/it-chain/it-chain-Engine/icode/domain/model"
	"github.com/it-chain/leveldb-wrapper"
)

type ICodeMetaRepository struct {
	leveldb *leveldbwrapper.DB
}

func NewICodeMetaRepository(path string) *ICodeMetaRepository {
	db := leveldbwrapper.CreateNewDB(path)

	return &ICodeMetaRepository{
		leveldb: db,
	}
}

func (i ICodeMetaRepository) Save(iCodeMeta model.ICodeMeta) error {
	return nil
}

func (i ICodeMetaRepository) Remove(id model.ICodeID) error {
	return nil
}

func (i ICodeMetaRepository) FindByID(id model.ICodeID) *model.ICodeMeta {
	return nil
}

func (i ICodeMetaRepository) FindAll() []model.ICodeMeta {
	return nil
}
