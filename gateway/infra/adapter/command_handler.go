package adapter

import (
	"log"

	"errors"

	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/it-chain-Engine/gateway/api"
	"github.com/it-chain/midgard"
)

var ErrInvalidCommand = errors.New("invalid command ")

type ConnectionApi interface {
	CreateConnection(address string) (gateway.Connection, error)
	CloseConnection(connectionID string) error
}

type ConnectionCommandHandler struct {
	connectionApi ConnectionApi
}

func (c ConnectionCommandHandler) HandleConnectionCreateCommand(command gateway.ConnectionCreateCommand) error {

	if !isValidCommand(command) {
		return ErrInvalidCommand
	}

	_, err := c.connectionApi.CreateConnection(command.Address)

	if err != nil {
		log.Printf("invalid address [%s]")
	}

	return nil
}

func (c ConnectionCommandHandler) HandleConnectionCloseCommand(command gateway.ConnectionCloseCommand) error {

	if !isValidCommand(command) {
		return ErrInvalidCommand
	}

	err := c.connectionApi.CloseConnection(command.GetID())

	if err != nil {
		log.Printf("fail to close connection: [%s]", err)
	}

	return nil
}

type GrpcCommandHandler struct {
	messageApi api.MessageApi
}

func (g GrpcCommandHandler) HandleGrpcDeliverCommand(command gateway.GrpcDeliverCommand) error {

	if !isValidCommand(command) {
		return ErrInvalidCommand
	}

	g.messageApi.DeliverMessage(command.Body, command.Protocol, command.Recipients...)

	return nil
}

func isValidCommand(command midgard.Command) bool {

	if command.GetID() == "" {
		log.Printf("invalid command id [%s]", command.GetID())
		return false
	}

	return true
}
