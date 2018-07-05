package txpool

import "github.com/it-chain/midgard"

type TxCreateCommand struct {
	midgard.CommandModel
	Jsonrpc string
	Method  TxDataType
	Params  Param
	ID      string
	ICodeID string
}

type ProposeBlockCommand struct {
	midgard.CommandModel
	Transactions []Transaction
}

type SendTransactionsCommand struct {
	midgard.CommandModel
	Transactions []*Transaction
	Leader
}

type GrpcDeliverCommand struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}
