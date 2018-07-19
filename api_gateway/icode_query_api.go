package api_gateway

import (
	"github.com/it-chain/engine/icode"
)

// this is an api only for querying current state which is repository of transaction
type ICodeQueryApi struct {
	transactionRepository ICodeMetaRepository
}

// this repository is a current state of all uncommitted transactions
type ICodeMetaRepository interface {
	FindAllMeta() []icode.Meta
	Save(transaction icode.Meta)
}

type ICodeEventHandler struct {
	repo ICodeMetaRepository
}
