package api

import (
	"log"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/it-chain-Engine/gateway"
)

type ConnectionApi struct {
	eventRepository gateway.EventRepository
	grpcService     gateway.GrpcService
}

func NewConnectionApi(eventRepository gateway.EventRepository, grpcService gateway.GrpcService) *ConnectionApi {
	return &ConnectionApi{
		eventRepository: eventRepository,
		grpcService:     grpcService,
	}
}

func (c ConnectionApi) CreateConnection(address string) (gateway.Connection, error) {

	log.Printf("dialing [%s]", address)

	connection, err := c.grpcService.Dial(address)

	if err != nil {
		log.Printf("fail to dial [%s]", err)
		return gateway.Connection{}, err
	}

	return gateway.NewConnection(connection.ID, connection.Address)
}

func (c ConnectionApi) CloseConnection(connectionID string) error {

	connection := &gateway.Connection{}

	err := eventstore.Load(connection, connectionID)

	if err != nil {
		return err
	}

	c.grpcService.CloseConnection(connectionID)

	return gateway.CloseConnection(connection.ID)
}

func (c ConnectionApi) OnConnection(connection gateway.Connection) (gateway.Connection, error) {

	return gateway.NewConnection(connection.ID, connection.Address)
}

func (c ConnectionApi) OnDisconnection(connection gateway.Connection) error {

	return gateway.CloseConnection(connection.ID)
}
