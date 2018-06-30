package gateway

import (
	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type Connection struct {
	midgard.AggregateModel
	Address string
}

func AddConnection(connection Connection) error {

	return eventstore.Save(connection.ID, ConnectionCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   connection.ID,
			Type: "connection.created",
		},
		Address: connection.Address,
	})
}

func CloseConnection(connectionID string) error {

	return eventstore.Save(connectionID, ConnectionClosedEvent{
		EventModel: midgard.EventModel{
			ID:   connectionID,
			Type: "connection.closed",
		},
	})
}
