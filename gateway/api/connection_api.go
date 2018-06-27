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

	events := make([]midgard.Event, 0)

	log.Printf("dialing [%s]", address)
	connection, err := c.grpcService.Dial(address)

	if err != nil {
		log.Printf("fail to dial [%s]", err)
		return err
	}

	events = append(events, gateway.ConnectionCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   connection.ID,
			Type: "connection.created",
		},
		Address: connection.Address,
	})

	err = c.eventRepository.Save(connection.ID, events...)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

//다른 node와의 연결 close
//todo close connection event 발생
func (c ConnectionApi) CloseConnection(connectionID string) error {

	events := make([]midgard.Event, 0)

	c.grpcService.CloseConnection(connectionID)

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

//
func (c ConnectionApi) OnConnection(connection gateway.Connection) error {

	events := make([]midgard.Event, 0)

	events = append(events, gateway.ConnectionCreatedEvent{
		EventModel: midgard.EventModel{
			ID: connection.ID,
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

//연결된 node의 connection 종료
//todo close connection event 발생
func (c ConnectionApi) OnDisconnection(connection gateway.Connection) error {

	events := make([]midgard.Event, 0)

	c.grpcService.CloseConnection(connection.GetID())

	events = append(events, gateway.ConnectionClosedEvent{
		EventModel: midgard.EventModel{
			ID:   connection.GetID(),
			Type: "connection.closed",
		},
	})

	err := c.eventRepository.Save(connection.GetID(), events...)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
