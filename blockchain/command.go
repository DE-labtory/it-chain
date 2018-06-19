package blockchain

import "github.com/it-chain/midgard"

type SyncUpdateCommand struct {
	midgard.EventModel
	sync bool
}



