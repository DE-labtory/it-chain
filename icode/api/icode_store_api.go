package api

import (
	"github.com/it-chain/it-chain-Engine/icode"
)

type ICodeStoreApi interface {
	//clone code from deploy info
	Clone(repositoryUrl string) (*icode.Meta, error)

	//push code to auth repo
	Push(meta icode.Meta) error
}
