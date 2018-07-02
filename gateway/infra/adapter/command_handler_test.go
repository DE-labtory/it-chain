package adapter_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/it-chain-Engine/gateway/infra/adapter"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
	"github.com/pkg/errors"
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

	apiErr := errors.New("fail to create connection")
	//given
	tests := map[string]struct {
		input    gateway.ConnectionCreateCommand
		mockFunc func(address string) (gateway.Connection, error)
		err      error
	}{
		"valid command": {
			input: gateway.ConnectionCreateCommand{
				CommandModel: midgard.CommandModel{
					ID: "conn1",
				},
				Address: "127.0.0.1:6666",
			},
			mockFunc: func(address string) (gateway.Connection, error) {
				assert.Equal(t, "127.0.0.1:6666", address)

				return gateway.Connection{}, nil
			},
			err: nil,
		},

		"invalid command": {
			input:    gateway.ConnectionCreateCommand{},
			mockFunc: nil,
			err:      adapter.ErrInvalidCommand,
		},

		"connectionApi error": {
			input: gateway.ConnectionCreateCommand{
				CommandModel: midgard.CommandModel{
					ID: "conn1",
				},
				Address: "127.0.0.1:6666",
			},
			mockFunc: func(address string) (gateway.Connection, error) {

				return gateway.Connection{}, apiErr
			},
			err: apiErr,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//given
		connectionApi := MockConnectionApi{}
		connectionApi.createConnectionFunc = test.mockFunc

		connectionCommandHandler := adapter.ConnectionCommandHandler{ConnectionApi: connectionApi}

		//when
		err := connectionCommandHandler.HandleConnectionCreateCommand(test.input)

		//then
		assert.Equal(t, err, test.err)
	}
}

func TestConnectionCommandHandler_HandleConnectionCreateCommand(t *testing.T) {

	apiErr := errors.New("fail to create connection")
	//given
	tests := map[string]struct {
		input    gateway.ConnectionCloseCommand
		mockFunc func(connectionID string) error
		err      error
	}{
		"valid command": {
			input: gateway.ConnectionCloseCommand{
				CommandModel: midgard.CommandModel{
					ID: "conn1",
				},
			},
			mockFunc: func(connectionID string) error {
				assert.Equal(t, "conn1", connectionID)

				return nil
			},
			err: nil,
		},

		"invalid command": {
			input:    gateway.ConnectionCloseCommand{},
			mockFunc: nil,
			err:      adapter.ErrInvalidCommand,
		},

		"connectionApi error": {
			input: gateway.ConnectionCloseCommand{
				CommandModel: midgard.CommandModel{
					ID: "conn1",
				},
			},
			mockFunc: func(connectionID string) error {

				return apiErr
			},
			err: apiErr,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//given
		connectionApi := MockConnectionApi{}
		connectionApi.closeConnectionFunc = test.mockFunc

		connectionCommandHandler := adapter.ConnectionCommandHandler{ConnectionApi: connectionApi}

		//when
		err := connectionCommandHandler.HandleConnectionCloseCommand(test.input)

		//then
		assert.Equal(t, err, test.err)
	}
}

type MockMessageApi struct {
	deliverMessageFunc func(body []byte, protocol string, ids ...string)
}

func (m MockMessageApi) DeliverMessage(body []byte, protocol string, ids ...string) {
	m.deliverMessageFunc(body, protocol, ids...)
}

func TestGrpcCommandHandler_HandleGrpcDeliverCommand(t *testing.T) {

	//given
	tests := map[string]struct {
		input    gateway.GrpcDeliverCommand
		mockFunc func(body []byte, protocol string, ids ...string)
		err      error
	}{
		"valid command": {
			input: gateway.GrpcDeliverCommand{
				CommandModel: midgard.CommandModel{
					ID: "conn1",
				},
			},
			mockFunc: func(body []byte, protocol string, ids ...string) {

			},
			err: nil,
		},

		"invalid command": {
			input:    gateway.GrpcDeliverCommand{},
			mockFunc: nil,
			err:      adapter.ErrInvalidCommand,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		//given
		messageApi := MockMessageApi{}
		messageApi.deliverMessageFunc = test.mockFunc

		connectionCommandHandler := adapter.GrpcCommandHandler{MessageApi: messageApi}

		//when
		err := connectionCommandHandler.HandleGrpcDeliverCommand(test.input)

		//then
		assert.Equal(t, err, test.err)

	}
}
