package blockchain

import (
	"time"

	"github.com/it-chain/midgard"
)

type SyncStartEvent struct {
	midgard.EventModel
}

type SyncDoneEvent struct {
	midgard.EventModel
}

// ToDo: 삭제 제안 - junk_sound
//type BlockAddToPoolEvent struct {
//	midgard.EventModel
//	Seal      []byte
//	PrevSeal  []byte
//	Height    uint64
//	TxList    []byte
//	TxSeal    [][]byte
//	Timestamp time.Time
//	Creator   []byte
//}
//
//type BlockRemoveFromPoolEvent struct {
//	midgard.EventModel
//	Height uint64
//}

// event when block is saved to event store
type BlockCommittedEvent struct {
	midgard.EventModel
	Seal string
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
}
