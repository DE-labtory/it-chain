package txpool

import "github.com/it-chain/midgard"

type TxCreateCommand struct {
	midgard.CommandModel
	TxData
	TxType TxDataType
}
