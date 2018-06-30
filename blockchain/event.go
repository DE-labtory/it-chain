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
}

type BlockValidatedEvent struct {
	midgard.EventModel
	Block Block
}
