package api

import (
	"log"

	"github.com/it-chain/it-chain-Engine/grpc_gateway"
	"github.com/it-chain/midgard"
)

type ConnectionApi struct {
	eventRepository midgard.Repository
	grpcService     grpc_gateway.GrpcService
}

func NewConnectionApi(eventRepository midgard.Repository, grpcService grpc_gateway.GrpcService) *ConnectionApi {
	return &ConnectionApi{
		eventRepository: eventRepository,
		grpcService:     grpcService,
	}
}

func (c ConnectionApi) CreateConnection(command grpc_gateway.ConnectionCreateCommand) {

	events := make([]midgard.Event, 0)

	if command.Address == "" {
		log.Printf("invalid address [%s]")
		return
	}

	log.Printf("dialing [%s]", command.Address)

	connection, err := c.grpcService.Dial(command.Address)

	if err != nil {
		log.Printf("fail to dial [%s]", err)
	}

	events = append(events, grpc_gateway.ConnectionCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   connection.ID,
			Type: "connection.created",
		},
		Address: connection.Address,
	})

	err = c.eventRepository.Save(connection.ID, events...)

	if err != nil {
		log.Println(err.Error())
	}
}

//다른 node와의 연결 close
//todo close connection event 발생
func (c ConnectionApi) CloseConnection(connection grpc_gateway.Connection) {

}

//
func (c ConnectionApi) OnConnection(connection grpc_gateway.Connection) {

	events := make([]midgard.Event, 0)

	events = append(events, grpc_gateway.ConnectionCreatedEvent{
		EventModel: midgard.EventModel{
			ID: connection.ID,
		},
		Address: connection.Address,
	})

	err := c.eventRepository.Save(connection.ID, events...)

	if err != nil {
		log.Println(err.Error())
	}
}

//연결된 node의 connection 종료
//todo close connection event 발생
func (c ConnectionApi) OnDisconnection(connection grpc_gateway.Connection) {

}
