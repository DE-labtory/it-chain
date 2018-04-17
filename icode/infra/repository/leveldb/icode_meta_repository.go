package leveldb

import (
	"encoding/json"

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

	b, err := iCodeMeta.Serialize()

	if err != nil {
		return err
	}

	if err = i.leveldb.Put([]byte(iCodeMeta.ID), b, true); err != nil {
		return err
	}

	return nil
}

func (i ICodeMetaRepository) Remove(id model.ICodeID) error {

	return i.leveldb.Delete([]byte(id), true)
}

func (i ICodeMetaRepository) FindByID(id model.ICodeID) (*model.ICodeMeta, error) {

	b, err := i.leveldb.Get([]byte(id))

	if err != nil {
		return nil, err
	}

	var iCodeMeta *model.ICodeMeta

	err = json.Unmarshal(b, iCodeMeta)

	if err != nil {
		return nil, err
	}

	return iCodeMeta, nil
}

func (i ICodeMetaRepository) FindAll() ([]*model.ICodeMeta, error) {

	iter := i.leveldb.GetIteratorWithPrefix([]byte(""))
	iCodeMetas := make([]*model.ICodeMeta, 0)

	for iter.Next() {
		val := iter.Value()
		iCodeMeta := &model.ICodeMeta{}
		err := deserialize(val, iCodeMeta)

		if err != nil {
			return nil, err
		}

		iCodeMetas = append(iCodeMetas, iCodeMeta)
	}

	return iCodeMetas, nil
}

func deserialize(b []byte, iCodeModel *model.ICodeMeta) error {
	err := json.Unmarshal(b, iCodeModel)

	if err != nil {
		return nil
	}
}
