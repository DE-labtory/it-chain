package blockchain

import "github.com/it-chain/midgard"

type SyncUpdateCommand struct {
	midgard.EventModel
	sync bool
}

type MessageDeliverCommand struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}
