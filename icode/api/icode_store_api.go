package api

import "github.com/it-chain/it-chain-Engine/icode/domain/icodeMeta"

//Api to import or store icode from outside
type ItCodeStoreApi interface {

	//get icode from outside
	Clone(repositoryUrl string) (*icodeMeta.ICode, error)

	//push code to auth repo
	Push(icodeMeta.ICode) error
}
