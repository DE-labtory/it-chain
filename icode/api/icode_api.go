package api

import "github.com/it-chain/it-chain-Engine/icode"

type ICodeApi struct {
	ContainerService icode.ContainerService
	StoreApi         ICodeStoreApi
	MetaRepository   icode.MetaRepository
}

func (iApi ICodeApi)()  {

}