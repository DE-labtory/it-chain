package blockchain

import "github.com/it-chain/midgard"

type SyncStartEvent struct {
	midgard.EventModel
}

type SyncDoneEvent struct {
	midgard.EventModel
}

type BlockAddToPoolEvent struct {
	midgard.EventModel
	BlockHeight uint64
}

type BlockRemoveFromPoolEvent struct {
	midgard.EventModel
	BlockHeight uint64
}
