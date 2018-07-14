package grpc_gateway

import "github.com/it-chain/midgard"

type ConnectionCreatedEvent struct {
	midgard.EventModel
	Address string
}

type ConnectionClosedEvent struct {
	midgard.EventModel
}

type ErrorCreatedEvent struct {
	midgard.EventModel
	Event string
	Err   string
}

type EventRepository interface {
	Load(aggregate midgard.Aggregate, aggregateID string) error
	Save(aggregateID string, events ...midgard.Event) error
}
