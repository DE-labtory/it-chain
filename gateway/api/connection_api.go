package api

import (
	"log"

	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/midgard"
)

type ConnectionApi struct {
	eventRepository midgard.Repository
	grpcService     gateway.GrpcService
}

func NewConnectionApi(eventRepository midgard.Repository, grpcService gateway.GrpcService) *ConnectionApi {
	return &ConnectionApi{
		eventRepository: eventRepository,
		grpcService:     grpcService,
	}
}

func (c ConnectionApi) CreateConnection(address string) error {

	log.Printf("dialing [%s]", address)

	connection, err := c.grpcService.Dial(address)

	if err != nil {
		log.Printf("fail to dial [%s]", err)
		return err
	}

	return c.storeConnectionCreatedEvent(connection)
}

func (c ConnectionApi) storeConnectionCreatedEvent(connection gateway.Connection) error {

	events := make([]midgard.Event, 0)

	events = append(events, gateway.ConnectionCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   connection.ID,
			Type: "connection.created",
		},
		Address: connection.Address,
	})

	err := c.eventRepository.Save(connection.ID, events...)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (c ConnectionApi) CloseConnection(connectionID string) error {

	c.grpcService.CloseConnection(connectionID)

	return c.storeConnectionClosedEvent(connectionID)
}

func (c ConnectionApi) storeConnectionClosedEvent(connectionID string) error {

	events := make([]midgard.Event, 0)

	events = append(events, gateway.ConnectionClosedEvent{
		EventModel: midgard.EventModel{
			ID:   connectionID,
			Type: "connection.closed",
		},
	})

	err := c.eventRepository.Save(connectionID, events...)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (c ConnectionApi) OnConnection(connection gateway.Connection) error {

	return c.storeConnectionCreatedEvent(connection)
}

func (c ConnectionApi) OnDisconnection(connection gateway.Connection) error {

	return c.storeConnectionClosedEvent(connection.ID)
}
