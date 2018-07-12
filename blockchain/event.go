package blockchain

import (
	"github.com/it-chain/midgard"
	"time"
)

type SyncStartEvent struct {
	midgard.EventModel
}

type SyncDoneEvent struct {
	midgard.EventModel
}

type BlockAddToPoolEvent struct {
	midgard.EventModel
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []byte
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   []byte
}

type BlockRemoveFromPoolEvent struct {
	midgard.EventModel
	Height uint64
}

// event when block is saved to event store
type BlockCommittedEvent struct {
	midgard.EventModel
	Seal string
}
