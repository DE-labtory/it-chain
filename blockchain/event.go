package blockchain

import "github.com/it-chain/midgard"

type NodeUpdateEvent struct {
	midgard.EventModel
}

type NodeCreatedEvent struct {
	midgard.EventModel
	Node
}

type NodeDeletedEvent struct {
	midgard.EventModel
	Node
}

type BlockValidatedEvent struct {
	midgard.EventModel
	Block Block
}
