package api

import (
	"log"

	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/engine/grpc_gateway"
)

type ConnectionApi struct {
	eventRepository grpc_gateway.EventRepository
	grpcService     grpc_gateway.GrpcService
}

func NewConnectionApi(eventRepository grpc_gateway.EventRepository, grpcService grpc_gateway.GrpcService) *ConnectionApi {
	return &ConnectionApi{
		eventRepository: eventRepository,
		grpcService:     grpcService,
	}
}

func (c ConnectionApi) CreateConnection(address string) (grpc_gateway.Connection, error) {

	log.Printf("dialing [%s]", address)

	connection, err := c.grpcService.Dial(address)

	if err != nil {
		log.Printf("fail to dial [%s]", err)
		return grpc_gateway.Connection{}, err
	}

	return grpc_gateway.NewConnection(connection.ID, connection.Address)
}

func (c ConnectionApi) CloseConnection(connectionID string) error {

	connection := &grpc_gateway.Connection{}

	err := eventstore.Load(connection, connectionID)

	if err != nil {
		return err
	}

	c.grpcService.CloseConnection(connectionID)

	return grpc_gateway.CloseConnection(connection.ID)
}

func (c ConnectionApi) OnConnection(connection grpc_gateway.Connection) (grpc_gateway.Connection, error) {

	return grpc_gateway.NewConnection(connection.ID, connection.Address)
}

func (c ConnectionApi) OnDisconnection(connection grpc_gateway.Connection) error {

	return grpc_gateway.CloseConnection(connection.ID)
}
