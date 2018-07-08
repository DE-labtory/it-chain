package service

import (
	"errors"
	"fmt"

	"github.com/it-chain/it-chain-Engine/icode"
	"github.com/it-chain/tesseract"
	"github.com/it-chain/tesseract/cellcode/cell"
)

type TesseractContainerService struct {
	tesseract      *tesseract.Tesseract
	repository     icode.ReadOnlyMetaRepository
	containerIdMap map[icode.ID]string // key : iCodeId, value : containerId
}

func NewTesseractContainerService(config tesseract.Config) *TesseractContainerService {
	tesseractObj := &TesseractContainerService{
		tesseract: tesseract.New(config),
	}
	tesseractObj.InitContainers()
	return tesseractObj
}

func (cs TesseractContainerService) StartContainer(meta icode.Meta) error {
	tesseractIcodeInfo := tesseract.ICodeInfo{
		Name:      meta.RepositoryName,
		Directory: meta.Path,
	}
	containerId, err := cs.tesseract.SetupContainer(tesseractIcodeInfo)
	if err != nil {
		return err
	}
	cs.containerIdMap[meta.ICodeID] = containerId
	return nil
}

func (cs TesseractContainerService) ExecuteTransaction(tx icode.Transaction) (*icode.Result, error) {
	containerId, found := cs.containerIdMap[tx.TxData.ICodeID]
	if !found {
		return nil, errors.New(fmt.Sprintf("no container for iCode : %s", tx.TxData.ICodeID))
	}
	tesseractTxInfo = cell.TxInfo{
		Method: tx.TxData.Method,
		ID:     tx.TxData.ID,
		Params: tx.TxData.Params,
	}
	cs.tesseract.QueryOrInvoke(containerId)
	return nil, nil
}

func (cs TesseractContainerService) StopContainer(id icode.ID) error {
	panic("implement please")
	return nil
}

// start containers in repos
func (cs *TesseractContainerService) InitContainers() error {
	metas, err := cs.repository.FindAll()
	if err != nil {
		return err
	}
	for _, meta := range metas {
		cs.StartContainer(*meta)
	}
}
