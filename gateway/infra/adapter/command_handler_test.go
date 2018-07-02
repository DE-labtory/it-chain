package adapter_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/it-chain-Engine/gateway/infra/adapter"
)

type MockConnectionApi struct {
	createConnectionFunc func(address string) (gateway.Connection, error)
	closeConnectionFunc  func(connectionID string) error
}

func (m MockConnectionApi) CreateConnection(address string) (gateway.Connection, error) {
	return m.createConnectionFunc(address)
}

func (m MockConnectionApi) CloseConnection(connectionID string) error {
	return m.closeConnectionFunc(connectionID)
}

func TestConnectionCommandHandler_HandleConnectionCloseCommand(t *testing.T) {

	//given
	tests := map[string]struct {
		input  gateway.ConnectionCreateCommand
		output gateway.GrpcReceiveCommand
		err    error
	}{
		"success": {
			input: gateway.ConnectionCreateCommand{},
			output: gateway.GrpcReceiveCommand{
				Body:         []byte("hello world"),
				ConnectionID: "123",
			},
			err: nil,
		},
	}

	connectionCommandHandler := adapter.ConnectionCommandHandler{}

	for testName, test := range tests {
		connectionCommandHandler.HandleConnectionCreateCommand()
	}

}

func TestConnectionCommandHandler_HandleConnectionCreateCommand(t *testing.T) {

}

func TestGrpcCommandHandler_HandleGrpcDeliverCommand(t *testing.T) {

}
