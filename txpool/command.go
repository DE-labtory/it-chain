package txpool

import "github.com/it-chain/midgard"

type TxCreateCommand struct {
	midgard.CommandModel
	TxData
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

type MessageDeliverCommand struct {
	midgard.CommandModel
	Transactions []*Transaction
	Leader
}