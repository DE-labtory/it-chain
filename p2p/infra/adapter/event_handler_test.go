package adapter_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/infra/adapter"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
)

type EventHandlerMockNodeApi struct {}
func (na EventHandlerMockNodeApi) AddNode(node p2p.Node) error{return nil}
func (na EventHandlerMockNodeApi) DeleteNode(id p2p.NodeId) error{return nil}


func TestEventHandler_HandleConnCreatedEvent(t *testing.T) {
	//1. test proper input
	//empty nodeid, empty address
	//2. proper output
	//matching err

	tests := map[string] struct{
		input struct{
			nodeId string
			address string
		}
		err error
	}{
		"success":{
			input: struct {
				nodeId  string
				address string
			}{nodeId: string("123"), address: string("123")},
			err:nil,
		},
		"empty nodeId test":{
			input: struct {
				nodeId  string
				address string
			}{nodeId: string(""), address: string("123")},
			err:adapter.ErrEmptyNodeId,
		},
		"empty address test":{
			input: struct {
				nodeId  string
				address string
			}{nodeId: string("123"), address: string("")},
			err:adapter.ErrEmptyAddress,
		},
	}
	eventHandler := adapter.NewEventHandler(EventHandlerMockNodeApi{})

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		err := eventHandler.HandleConnCreatedEvent(p2p.ConnectionCreatedEvent{EventModel:midgard.EventModel{ID:test.input.nodeId}, Address:test.input.address})
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
			err: adapter.ErrEmptyNodeId,
		},
	}

	eventHandler := adapter.NewEventHandler(EventHandlerMockNodeApi{})
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

type MockNodeRepository struct{}
func (nr MockNodeRepository) Save(data p2p.Node) error{return nil}
type MockLeaderRepository struct{}
func (lr MockLeaderRepository) SetLeader(leader p2p.Leader){}


func TestRepositoryProjector_HandleLeaderUpdatedEvent(t *testing.T) {
	repositoryProjector := SetupRepositoryProjector()

	tests := map[string]struct {
		input struct {
			id string
		}
		err error
	}{
		"success": {
			input: struct {
				id string
			}{id: "123"},
			err: nil,
		},
		"empty node id test": {
			input: struct {
				id string
			}{id: string("")},
			err: adapter.ErrEmptyNodeId,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		event := p2p.LeaderUpdatedEvent{
			EventModel: midgard.EventModel{
				ID: test.input.id,
			},
		}
		err := repositoryProjector.HandleLeaderUpdatedEvent(event)
		assert.Equal(t, err, test.err)
	}
}

func SetupRepositoryProjector() (*adapter.RepositoryProjector) {

	repositoryProjector := adapter.NewRepositoryProjector(MockNodeRepository{}, MockLeaderRepository{})

	return repositoryProjector
}
