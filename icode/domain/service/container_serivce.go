package service

import "github.com/it-chain/it-chain-Engine/icode/domain/model"

type ContainerService interface {
	Start(icodeMeta model.ICodeMeta) error
	Stop(id model.ICodeID) error
	Run(tx model.Transaction) error
}
