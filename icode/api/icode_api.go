package api

import "github.com/it-chain/it-chain-Engine/icode"

type ICodeApi struct {
	ContainerService icode.ContainerService
	StoreApi         ICodeStoreApi
	MetaRepository   icode.MetaRepository
}

func (iApi ICodeApi) Deploy(gitUrl string) error {
	panic("implement please")
}

func (iApi ICodeApi) UnDeploy(id icode.ID) error {
	panic("implement please")
}

func (iApi ICodeApi) Invoke(txs []icode.Transaction) {
	panic("implement please")
}

func (iApi ICodeApi) Query(tx icode.Transaction) (icode.Result, error) {
	panic("implement please")
}
