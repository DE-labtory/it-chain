package leveldb

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/icode"
	"github.com/it-chain/leveldb-wrapper"
)

func NewMetaRepository(path string) *MetaRepository {
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return &MetaRepository{
		levelDb: db,
	}
}

type MetaRepository struct {
	levelDb *leveldbwrapper.DB
}

func (mr *MetaRepository) Save(meta icode.Meta) error {

	if meta.ICodeID == "" {
		return errors.New("ICodeMeta ID is empty")
	}

	b, err := common.Serialize(meta)

	if err != nil {
		return err
	}

	if err = mr.levelDb.Put([]byte(meta.ICodeID), b, true); err != nil {
		return err
	}

	return nil
}

func (mr *MetaRepository) Remove(id icode.ID) error {
	return mr.levelDb.Delete([]byte(id), true)
}

func (mr *MetaRepository) FindById(id icode.ID) (*icode.Meta, error) {
	b, err := mr.levelDb.Get([]byte(id))

	if err != nil {
		return nil, errors.New("NO")
	}

	if len(b) == 0 {
		return nil, nil
	}

	iCodeMeta := &icode.Meta{}

	err = common.Deserialize(b, iCodeMeta)

	if err != nil {
		return nil, err
	}

	return iCodeMeta, nil
}

func (mr *MetaRepository) FindAll() ([]*icode.Meta, error) {
	iter := mr.levelDb.GetIteratorWithPrefix([]byte(""))
	iCodeMetas := make([]*icode.Meta, 0)

	for iter.Next() {
		val := iter.Value()
		iCodeMeta := &icode.Meta{}
		err := common.Deserialize(val, iCodeMeta)

		if err != nil {
			return nil, err
		}

		iCodeMetas = append(iCodeMetas, iCodeMeta)
	}

	return iCodeMetas, nil
}

func (mr *MetaRepository) FindByGitURL(url string) (*icode.Meta, error) {
	metas, err := mr.FindAll()
	if err != nil {
		return nil, err
	}

	for _, meta := range metas {
		if meta.GitUrl == url {
			return meta, nil
		}
	}

	return &icode.Meta{}, nil
}
