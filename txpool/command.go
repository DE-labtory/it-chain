package txpool

import "github.com/it-chain/midgard"

type TxCreateCommand struct {
	midgard.CommandModel
	TxData
}

func (t TxCreateCommand) GetID() string {
	return t.TxData.ID
}
