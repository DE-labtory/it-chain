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

type BlockValidatedEvent struct {
	midgard.EventModel
	Block Block
}

type BlockProposedEvent struct {
	midgard.EventModel
	Block Block
}
