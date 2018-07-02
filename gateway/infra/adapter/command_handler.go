package adapter

import (
	"log"

	"errors"

	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/midgard"
)

var ErrInvalidCommand = errors.New("invalid command ")

type ConnectionApi interface {
	CreateConnection(address string) (gateway.Connection, error)
	CloseConnection(connectionID string) error
}

type ConnectionCommandHandler struct {
	ConnectionApi ConnectionApi
}

func (c ConnectionCommandHandler) HandleConnectionCreateCommand(command gateway.ConnectionCreateCommand) error {

	if !isValidCommand(command) {
		return ErrInvalidCommand
	}

	_, err := c.ConnectionApi.CreateConnection(command.Address)

	if err != nil {
		log.Printf("invalid address [%s]")
		return err
	}

	return nil
}

func (c ConnectionCommandHandler) HandleConnectionCloseCommand(command gateway.ConnectionCloseCommand) error {

	if !isValidCommand(command) {
		return ErrInvalidCommand
	}

	err := c.ConnectionApi.CloseConnection(command.GetID())

	if err != nil {
		log.Printf("fail to close connection: [%s]", err)
		return err
	}

	return nil
}

type MessageApi interface {
	DeliverMessage(body []byte, protocol string, ids ...string)
}

type GrpcCommandHandler struct {
	MessageApi MessageApi
}

func (g GrpcCommandHandler) HandleGrpcDeliverCommand(command gateway.GrpcDeliverCommand) error {

	if !isValidCommand(command) {
		return ErrInvalidCommand
	}

	g.MessageApi.DeliverMessage(command.Body, command.Protocol, command.Recipients...)

	return nil
}

func isValidCommand(command midgard.Command) bool {

	if command.GetID() == "" {
		log.Printf("invalid command id [%s]", command.GetID())
		return false
	}

	return true
}
