package p2p

import "github.com/it-chain/midgard"

//publish
type LeaderChangeEvent struct {
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
}

// handle leader received event
type LeaderUpdatedEvent struct {
	midgard.EventModel
	Leader Leader
}

//handle
//node list received
type NodeListUpdatedEvent struct {
	midgard.EventModel
	NodeList []Node
}