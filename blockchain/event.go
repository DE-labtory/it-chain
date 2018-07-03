package blockchain

import "github.com/it-chain/midgard"

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

type SyncStartEvent struct {
	midgard.EventModel
}

type SyncDoneEvent struct {
	midgard.EventModel
}
