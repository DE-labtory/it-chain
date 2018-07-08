package icode

import "github.com/it-chain/midgard"

type TransactionExecuteCommand struct {
	midgard.CommandModel
	transaction Transaction
}
