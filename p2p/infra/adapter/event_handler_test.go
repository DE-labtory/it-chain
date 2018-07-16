package adapter_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/infra/adapter"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
	"github.com/it-chain/it-chain-Engine/p2p/test/mock"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
)

func TestEventHandler_HandleConnCreatedEvent(t *testing.T) {

	tests := map[string]struct {
		input struct {
			nodeId  string
			address string
		}
		err error
	}{
		"success": {
			input: struct {
				nodeId  string
				address string
			}{nodeId: string("123"), address: string("123")},
			err: eventstore.ErrNilStore,
		},
		"empty address test": {
			input: struct {
				nodeId  string
				address string
			}{nodeId: string("123"), address: string("")},
			err: p2p.ErrEmptyAddress,
		},
	}

	communicationApi := &mock.MockCommunicationApi{}

	communicationApi.DeliverPLTableFunc = func(connectionId string) error {
		return nil
	}

	eventHandler := adapter.NewEventHandler(communicationApi)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := eventHandler.HandleConnCreatedEvent(p2p.ConnectionCreatedEvent{EventModel: midgard.EventModel{ID: test.input.nodeId}, Address: test.input.address})
		assert.Equal(t, err, test.err)
	}

}

func TestEventHandler_HandleConnDisconnectedEvent(t *testing.T) {

	tests := map[string]struct {
		input struct {
			id string
		}
		err error
	}{
		"success": {
			input: struct {
				id string
			}{id: string(123)},
			err: nil,
		},
		"empty node id test": {
			input: struct {
				id string
			}{id: string("")},
			err: adapter.ErrEmptyPeerId,
		},
	}

	communicationApi := &mock.MockCommunicationApi{}

	eventHandler := adapter.NewEventHandler(communicationApi)

	for testName, test := range tests {

		t.Logf("running test case %s", testName)

		event := p2p.ConnectionDisconnectedEvent{
			EventModel: midgard.EventModel{
				ID: test.input.id,
			},
		}

		err := eventHandler.HandleConnDisconnectedEvent(event)

		assert.Equal(t, err, test.err)

	}
}
