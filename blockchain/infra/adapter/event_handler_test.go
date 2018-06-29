package adapter_test

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/infra/adapter"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/it-chain/midgard"
	"errors"
)

type MockNodeRepository struct {
	AddNodeFunc func(node blockchain.Node) error
	DeleteNodeFunc func(node blockchain.Node) error
}

func (nr MockNodeRepository) AddNode(node blockchain.Node) error {
	return nr.AddNodeFunc(node)
}
func (nr MockNodeRepository) DeleteNode(node blockchain.Node) error {
	return nr.DeleteNodeFunc(node)
}


func TestEventHandler_HandleNodeCreatedEvent(t *testing.T) {
	tests := map[string] struct {
		input struct {
			ID string
			nodeId string
			address string
			rpErr error
		}
		err error
	}{
		"success": {
			input: struct {
				ID string
				nodeId string
				address string
				rpErr error
			}{ID: string("zf"), nodeId: string("zf2"), address: string("11.22.33.44"), rpErr: nil},
			err: nil,
		},
		"empty eventId test": {
			input: struct {
				ID string
				nodeId string
				address string
				rpErr error
			}{ID: string(""), nodeId: string("zf2"), address: string("11.22.33.44"), rpErr: nil},
			err: adapter.ErrEmptyEventId,
		},
		"repository error test": {
			input: struct {
				ID string
				nodeId string
				address string
				rpErr error
			}{ID: string("zf"), nodeId: string("zf2"), address: string("11.22.33.44"), rpErr: errors.New("repository error")},
			err: adapter.ErrNodeApi,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		nr := MockNodeRepository{}
		nr.AddNodeFunc = func(node blockchain.Node) error {
			assert.Equal(t, node.NodeId.Id, string("zf2"))
			assert.Equal(t, node.IpAddress, string("11.22.33.44"))
			return test.input.rpErr
		}
		rp := adapter.RepositoryProjector{
			NodeRepository: nr,
		}

		eventHandler := adapter.NewEventHandler(rp)

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

func TestEventHandler_HandleNodeDeletedEvent(t *testing.T) {
	tests := map[string] struct {
		input struct {
			ID string
			nodeId string
			address string
			rpErr error
		}
		err error
	}{
		"success": {
			input: struct {
				ID string
				nodeId string
				address string
				rpErr error
			}{ID: string("zf"), nodeId: string("zf2"), address: string("11.22.33.44"), rpErr: nil},
			err: nil,
		},
		"empty eventId test": {
			input: struct {
				ID string
				nodeId string
				address string
				rpErr error
			}{ID: string(""), nodeId: string("zf2"), address: string("11.22.33.44"), rpErr: nil},
			err: adapter.ErrEmptyEventId,
		},
		"repository error test": {
			input: struct {
				ID string
				nodeId string
				address string
				rpErr error
			}{ID: string("zf"), nodeId: string("zf2"), address: string("11.22.33.44"), rpErr: errors.New("repository error")},
			err: adapter.ErrNodeApi,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		nr := MockNodeRepository{}
		nr.DeleteNodeFunc = func(node blockchain.Node) error {
			assert.Equal(t, node.NodeId.Id, string("zf2"))
			assert.Equal(t, node.IpAddress, string("11.22.33.44"))
			return test.input.rpErr
		}
		rp := adapter.RepositoryProjector{
			NodeRepository: nr,
		}

		eventHandler := adapter.NewEventHandler(rp)

		event := blockchain.NodeDeletedEvent{
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
		err := eventHandler.HandleNodeDeletedEvent(event)

		assert.Equal(t, err, test.err)
	}
}