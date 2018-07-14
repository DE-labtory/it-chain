package icode

import (
	"github.com/it-chain/midgard"
)

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

type BlockExecuteCommand struct {
	midgard.CommandModel
	Block []byte
}

type BlockResultCommand struct {
	midgard.CommandModel
	TxResults []Result
}
