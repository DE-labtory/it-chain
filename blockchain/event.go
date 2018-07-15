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

type BlockRemoveFromPoolEvent struct {
	midgard.EventModel
	Height uint64
}

type BlockStagedEvent struct {
	midgard.EventModel
	State string
}

type BlockCommittedEvent struct {
	midgard.EventModel
	State string
}

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
