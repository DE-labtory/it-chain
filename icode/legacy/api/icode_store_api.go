package api

import "github.com/it-chain/it-chain-Engine/icode/domain/model"

//Api to import or store icode from outside
type ItCodeStoreApi interface {
	Clone(repositoryUrl string) (*model.ICodeMeta, error)

	//push code to auth repo
	Push(meta model.ICodeMeta) error
}
