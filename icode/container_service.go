package icode

import "github.com/it-chain/it-chain-Engine/icode/legacy/domain/model"

type ContainerService interface {
	Start(meta Meta) error
	Stop(id ID) error
	Run(tx Transaction) (*model.Result, error)
}
