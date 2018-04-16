package api

import (
	"github.com/it-chain/it-chain-Engine/icode/domain/model"
	"github.com/it-chain/it-chain-Engine/icode/domain/repository"
)

type ICodeApi struct {
	ICodeMetaRepository repository.ICodeMetaRepository
}

func Deploy(ICodeMeta model.ICodeMeta) error {

	return nil
}

func ProposalTransactions(txs []model.Transaction) {

}

func QueryICode(requests model.Transaction) {

}

func InvokeICode(tx model.Transaction) {

}
