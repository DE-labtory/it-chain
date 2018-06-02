package api

import (
	"log"

	"github.com/it-chain/bifrost/client"
	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/midgard"
)

type ConnectionApi struct {
	eventRepository midgard.Repository
	dialService     gateway.DialService
}

func NewConnectionApi(eventRepository midgard.Repository, dialService gateway.DialService) *ConnectionApi {
	return &ConnectionApi{
		eventRepository: eventRepository,
		dialService:     dialService,
	}
}

// 새로운 connection 이 생성되면 처리하는 함수이다.
func (c ConnectionApi) CreateConnection(command gateway.ConnectionCreateCommand) {

	if command.Address == "" {
		return
	}

	log.Printf("dialing [%s]", command.Address)

	clientOpt := client.ClientOpts{
		Ip:     command.Address,
		PriKey: c.priKey,
		PubKey: c.pubKey,
	}

	grpcOpt := client.GrpcOpts{
		TlsEnabled: false,
		Creds:      nil,
	}

	connection, err := client.Dial(command.Address, clientOpt, grpcOpt)

	if err != nil {

		c.publisher.Publish("Event", "Error", gateway.ErrorCreatedEvent{
			Err:   err.Error(),
			Event: "Connection fail to create",
		})

		return
	}

	if c.store.GetConnection(connection.GetID()) != nil {
		log.Printf("same connection is existed")
		return
	}

	err = c.publisher.Publish("Event", "Connection", gateway.ConnectionCreatedEvent{
		Address: connection.GetIP(),
		EventModel: midgard.EventModel{
			ID: connection.GetID(),
		},
	})

	if err != nil {
		log.Println(err.Error())
		return
	}

	connection.Handle(NewRequestHandler(c.publisher))
	c.store.AddConnection(connection)

	go func() {
		defer connection.Close()

		if err := connection.Start(); err != nil {
			connection.Close()
		}

		log.Printf("connections are closing")
		c.store.DeleteConnection(connection.GetID())
	}()
}
