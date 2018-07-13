package blockchain

import (
	"time"

	"github.com/it-chain/midgard"
)

type NodeUpdateEvent struct {
	midgard.EventModel
}

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

type BlockCreatedEvent struct {
	midgard.EventModel
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []byte
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   []byte
	State     string
}

type BlockStagedEvent struct {
	midgard.EventModel
	State string
}

type BlockCommittedEvent struct {
	midgard.EventModel
	State string
}
