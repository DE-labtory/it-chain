package adapter_test

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/infra/adapter"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/it-chain/midgard"
)

type EventHandlerMockNodeApi struct {
	AddNodeFunc func(node blockchain.Node) error
}
func (na EventHandlerMockNodeApi) AddNode(node blockchain.Node) error {
	return nil
}
func (na EventHandlerMockNodeApi) DeleteNode(node blockchain.Node) error {
	return nil
}

func TestEventHandler_HandleNodeCreatedEvent(t *testing.T) {
	tests := map[string] struct {
		input struct {
			ID string
			nodeId string
			address string
		}
		err error
	}{
		"success": {
			input: struct {
				ID string
				nodeId string
				address string
			}{ID: string("zf"), nodeId: string("zf2"), address: string("11.22.33.44")},
			err: nil,
		},
		"empty eventId test": {
			input: struct {
				ID string
				nodeId string
				address string
			}{ID: string(""), nodeId: string("zf2"), address: string("11.22.33.44")},
			err: adapter.ErrEmptyEventId,
		},
	}

	mockNodeApi := EventHandlerMockNodeApi{}
	mockNodeApi.AddNodeFunc = func(node blockchain.Node) error {
		assert.Equal(t, node.NodeId.Id, string("zf2"))
		assert.Equal(t, node.IpAddress, string("11.22.33.44"))
		return nil
	}

	eventHandler := adapter.NewEventHandler(mockNodeApi)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		event := blockchain.NodeCreatedEvent{
			EventModel: midgard.EventModel{
				ID: test.input.ID,
			},
			Node: blockchain.Node{
				NodeId: blockchain.NodeId{
					test.input.nodeId,
				},
				IpAddress: test.input.address,
			},
		}
		err := eventHandler.HandleNodeCreatedEvent(event)

		assert.Equal(t, err, test.err)
	}
}