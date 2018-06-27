package adapter

import (
	"log"

	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/it-chain-Engine/gateway/api"
)

type ConnectionCommandHandler struct {
	connectionApi api.ConnectionApi
}

func (c ConnectionCommandHandler) HandleConnectionCreateCommand(command gateway.ConnectionCreateCommand) {

	if command.Address == "" {
		log.Printf("invalid address [%s]")
		return
	}

	err := c.connectionApi.CreateConnection(command.Address)

	if err != nil {
		log.Printf("invalid address [%s]")
	}
}

func (c ConnectionCommandHandler) HandleConnectionCloseCommand(command gateway.ConnectionCloseCommand) {

	if command.GetID() == "" {
		log.Printf("invalid connection id [%s]", command.GetID())
		return
	}

	err := c.connectionApi.CloseConnection(command.GetID())

	if err != nil {
		log.Printf("fail to close connection: [%s]", err)
	}
}

type MessageCommandHandler struct {
	messageApi api.MessageApi
}

func (c ConnectionCommandHandler) HandleGrpcDeliverCommand(command gateway.GrpcDeliverCommand) {

}
