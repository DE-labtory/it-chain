package service

import "github.com/it-chain/it-chain-Engine/icode/domain/model"

type TesseractContainerService struct {
}

func (cs TesseractContainerService) Start(icodeMeta model.ICodeMeta) error {

	return nil
}

func (cs TesseractContainerService) Run(tx model.Transaction) (*model.Result, error) {
	return nil, nil
}

func (cs TesseractContainerService) Stop(id model.ICodeID) error {
	return nil
}
