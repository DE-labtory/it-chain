package icode

import "github.com/it-chain/midgard"

type TransactionExecuteCommand struct {
	midgard.CommandModel
	transaction Transaction
}

type DeployCommand struct {
	midgard.CommandModel
	Url string
}
type UnDeployCommand struct {
	midgard.CommandModel
}
