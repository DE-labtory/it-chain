package service

import (
	"github.com/it-chain/it-chain-Engine/icode"
)

type TesseractContainerService struct {
}

func (cs TesseractContainerService) StartContainer(meta icode.Meta) error {
	panic("implement please")
	return nil
}

func (cs TesseractContainerService) ExecuteTransaction(tx icode.Transaction) (*icode.Result, error) {
	panic("implement please")
	return nil, nil
}

func (cs TesseractContainerService) StopContainer(id icode.ID) error {
	panic("implement please")
	return nil
}
