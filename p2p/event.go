package p2p

import "github.com/it-chain/midgard"

//publish
type LeaderChangedEvent struct {
	midgard.EventModel
}

//handle
type ConnectionCreatedEvent struct {
	midgard.EventModel
	Address string
}

//handle
type ConnectionDisconnectedEvent struct {
	midgard.EventModel
}

// node created event
type NodeCreatedEvent struct {
	midgard.EventModel
	IpAddress string
}

type NodeDeletedEvent struct {
	midgard.EventModel
}

// handle leader received event
type LeaderUpdatedEvent struct {
	midgard.EventModel
}
