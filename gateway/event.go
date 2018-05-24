package gateway

import "github.com/it-chain/midgard"

type ConnectionCreatedEvent struct {
	midgard.EventModel
	Address string
}

type ConnectionDisconnectedEvent struct {
	midgard.EventModel
}

type ErrorCreatedEvent struct {
	midgard.EventModel
	Event string
	Err   string
}
