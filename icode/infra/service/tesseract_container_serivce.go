package service

import (
	"github.com/it-chain/it-chain-Engine/icode"
)

type TesseractContainerService struct {
}

func (cs TesseractContainerService) Start(meta icode.Meta) error {
	panic("implement please")
	return nil
}

func (cs TesseractContainerService) Run(tx icode.Transaction) (*icode.Result, error) {
	panic("implement please")
	return nil, nil
}

func (cs TesseractContainerService) Stop(id icode.ID) error {
	panic("implement please")
	return nil
}
