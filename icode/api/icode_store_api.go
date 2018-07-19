package api

import (
	"github.com/it-chain/engine/icode"
)

type ICodeStoreApi interface {
	//clone code from deploy info
	Clone(id string, baseSavePath string, repositoryUrl string, sshPath string) (*icode.Meta, error)

	//push code to auth repo
	Push(meta icode.Meta) error
}
