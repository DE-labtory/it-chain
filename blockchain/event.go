package blockchain

import "github.com/it-chain/midgard"

// synchronize를 시작하거나 끝났을 때 event ID로 사용한다.
var BC_SYNC_STATE_EID = "BC_SYNC_STATE_EID"

// block queued event id
var BLOCK_QUEUED_EID = "BLOCK_QUEUED_EID"

type NodeUpdateEvent struct {
	midgard.EventModel
}

type NodeCreatedEvent struct {
	midgard.EventModel
	Peer
}

type NodeDeletedEvent struct {
	midgard.EventModel
	Peer
}

type BlockQueuedEvent struct {
	midgard.EventModel
	Block
}

type BlockValidatedEvent struct {
	midgard.EventModel
	Block Block
}
