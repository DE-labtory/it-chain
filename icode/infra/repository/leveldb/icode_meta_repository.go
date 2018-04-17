package leveldb

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/icode/domain/model"
	"github.com/it-chain/leveldb-wrapper"
	"github.com/pkg/errors"
)

type ICodeMetaRepository struct {
	leveldb *leveldbwrapper.DB
}

func NewICodeMetaRepository(path string) *ICodeMetaRepository {

	db := leveldbwrapper.CreateNewDB(path)
	db.Open()

	return &ICodeMetaRepository{
		leveldb: db,
	}
}

func (i ICodeMetaRepository) Save(iCodeMeta model.ICodeMeta) error {

	if iCodeMeta.ID.ToString() == "" {
		return errors.New("ICodeMeta ID is empty")
	}

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
		return nil, errors.New("NO")
	}

	if len(b) == 0 {
		return nil, nil
	}

	iCodeMeta := &model.ICodeMeta{}

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
		err := Deserialize(val, iCodeMeta)

		if err != nil {
			return nil, err
		}

		iCodeMetas = append(iCodeMetas, iCodeMeta)
	}

	return iCodeMetas, nil
}

func (i ICodeMetaRepository) Close() {
	i.leveldb.Close()
}

func Deserialize(b []byte, iCodeModel *model.ICodeMeta) error {
	err := json.Unmarshal(b, iCodeModel)

	if err != nil {
		return err
	}

	return nil
}
