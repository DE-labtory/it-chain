package repository

import "github.com/it-chain/it-chain-Engine/icode/domain/model"

type ICodeMetaRepository interface {
	Save(iCodeMeta model.ICodeMeta) error
	Remove(id model.ICodeID) error
	FindByID(id model.ICodeID) (*model.ICodeMeta, error)
	FindAll() ([]*model.ICodeMeta, error)
}
