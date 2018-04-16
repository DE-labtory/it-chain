package api

import (
	"github.com/it-chain/it-chain-Engine/icode/domain/model"
	"github.com/it-chain/it-chain-Engine/icode/domain/repository"
)

//ICode의 Invoke, Query, 검증 수행
type ICodeApi struct {
	ICodeMetaRepository repository.ICodeMetaRepository
}

//Deploy ICode from git and push to backup server
func Deploy(ICodeMeta model.ICodeMeta) error {

	return nil
}

//Invoke transactions on ICode
func Invoke(txs []model.Transaction) {

}

//Query transactions on ICode (Read Only transaction request on ICode)
func Query(tx model.Transaction) {

}
