package service

import "github.com/it-chain/it-chain-Engine/icode/domain/model"

type ContainerService interface {
	Start(icodeMeta model.ICodeMeta)
	Run(tx model.Transaction)
}
