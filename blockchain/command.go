package blockchain

import (
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/midgard"
)

type SyncUpdateCommand struct {
	midgard.EventModel
	sync bool
}

type NodeUpdateCommand struct {
	midgard.EventModel
	Node
}

type ProposeBlockCommand struct {
	midgard.CommandModel
	// TODO: Transaction이 너무 다름.
	Transactions []txpool.Transaction
}

type BlockValidateCommand struct {
	midgard.CommandModel
	Block Block
}

type GrpcCommand struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}
