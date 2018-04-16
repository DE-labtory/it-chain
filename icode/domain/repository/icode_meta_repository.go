package repository

import "github.com/it-chain/it-chain-Engine/icode/domain/model"

type ICodeMetaRepository interface {
	Save(iCodeMeta model.ICodeMeta) error
	FindByID(id model.ICodeID)
	FindAll() []model.ICodeMeta
}
