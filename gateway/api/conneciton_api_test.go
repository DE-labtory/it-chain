package api_test

import (
	"testing"

	"os"

	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/it-chain-Engine/gateway/api"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

type MockEventRepository struct {
	loadFunc func(aggregate midgard.Aggregate, aggregateID string) error
	saveFunc func(aggregateID string, events ...midgard.Event) error
}

func (m MockEventRepository) Load(aggregate midgard.Aggregate, aggregateID string) error {
	return m.loadFunc(aggregate, aggregateID)
}

func (m MockEventRepository) Save(aggregateID string, events ...midgard.Event) error {
	return m.saveFunc(aggregateID, events...)
}

type MockGrpcService struct {
	dialFunc            func(address string) (gateway.Connection, error)
	closeConnectionFunc func(connID string)
	sendMessagesFunc    func(message []byte, protocol string, connIDs ...string)
}

func (m MockGrpcService) Dial(address string) (gateway.Connection, error) {
	return m.dialFunc(address)
}

func (m MockGrpcService) CloseConnection(connID string) {
	m.closeConnectionFunc(connID)
}

func (m MockGrpcService) SendMessages(message []byte, protocol string, connIDs ...string) {
	m.sendMessagesFunc(message, protocol, connIDs...)
}

func init() {
}

func TestConnectionApi_CreateConnection(t *testing.T) {

	defer InitStore()()

	//given
	tests := map[string]struct {
		input string
		err   error
	}{
		"create connection success": {
			input: "127.0.0.1:8888",
			err:   nil,
		},
	}

	repo := MockEventRepository{}
	repo.saveFunc = func(aggregateID string, events ...midgard.Event) error {
		assert.Equal(t, aggregateID, "randomID")
		assert.Equal(t, len(events), 1)
		assert.IsType(t, events[0], gateway.ConnectionCreatedEvent{})
		address := events[0].(gateway.ConnectionCreatedEvent).Address
		assert.Equal(t, address, "127.0.0.1:8888")

		return nil
	}

	grpcSerivce := MockGrpcService{}
	grpcSerivce.dialFunc = func(address string) (gateway.Connection, error) {

		return gateway.Connection{
			AggregateModel: midgard.AggregateModel{
				ID: "randomID",
			},
			Address: address,
		}, nil
	}

	connectionApi := api.NewConnectionApi(repo, grpcSerivce)

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		_, err := connectionApi.CreateConnection(test.input)
		assert.Equal(t, err, test.err)
	}
}

func TestConnectionApi_CloseConnection(t *testing.T) {

	defer InitStore()()

	//given
	tests := map[string]struct {
		input string
		err   error
	}{
		"close connection success": {
			input: "conn1",
			err:   nil,
		},
	}

	repo := MockEventRepository{}
	repo.saveFunc = func(aggregateID string, events ...midgard.Event) error {

		assert.Equal(t, aggregateID, "conn1")
		assert.Equal(t, len(events), 1)
		assert.IsType(t, events[0], gateway.ConnectionClosedEvent{})

		return nil
	}

	grpcSerivce := MockGrpcService{}
	grpcSerivce.closeConnectionFunc = func(connID string) {
		assert.Equal(t, connID, "conn1")
	}

	connectionApi := api.NewConnectionApi(repo, grpcSerivce)

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)
		eventstore.Save(test.input, gateway.ConnectionCreatedEvent{
			EventModel: midgard.EventModel{
				ID:   test.input,
				Type: "connection.created",
			},
			Address: "123",
		})
		err := connectionApi.CloseConnection(test.input)
		assert.Equal(t, err, test.err)
	}
}

func InitStore() func() {

	path := "./.test"

	eventstore.InitLevelDBStore(path, nil,
		gateway.ConnectionCreatedEvent{},
		gateway.ConnectionClosedEvent{})

	return func() {
		eventstore.Close()
		os.RemoveAll(path)
	}
}
