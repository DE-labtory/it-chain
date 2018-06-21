package blockchain

import "github.com/it-chain/midgard"

type SyncUpdateCommand struct {
	midgard.EventModel
	sync bool
}

type NodeUpdateCommand struct {
	midgard.EventModel
}

