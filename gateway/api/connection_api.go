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

func (c ConnectionApi) CreateConnection(command gateway.ConnectionCreateCommand) {

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

	events = append(events, gateway.ConnectionCreatedEvent{
		EventModel: midgard.EventModel{
			ID: connection.ID,
		},
		Address: connection.Address,
	})

	err = c.eventRepository.Save(connection.ID, events...)

	if err != nil {
		log.Println(err.Error())
	}
}
