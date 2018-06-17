package adapter_test

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/p2p/infra/adapter"
	"github.com/it-chain/it-chain-Engine/p2p/infra/repository/leveldb"
	"github.com/it-chain/midgard"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/magiconair/properties/assert"
)



func TestRepositoryProjector_HandleConnCreatedEvent(t *testing.T) {
	repositoryProjector, endUp := SetupRepositoryProjector("path_node_repository", "path_leader_repository")

	tests := map[string] struct{
		input struct{
			id string
			address string
		}
		err error
	}{
		"success":{
			input: struct {
				id      string
				address string
			}{id: string(123), address: string(123)},
			err: nil,
		},
		"empty node id test":{
			input: struct {
				id      string
				address string
			}{id: string(""), address: string(123)},
		},
		"empty address test":{
			input: struct {
				id      string
				address string
			}{id: string(123), address: string("")},
		},
	}

	defer endUp()

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		event := p2p.ConnectionCreatedEvent{
			EventModel: midgard.EventModel{
				ID:test.input.id,
			},
			Address:test.input.address,
		}
		err := repositoryProjector.HandleConnCreatedEvent(event)
		assert.Equal(t, err, test.err)

		node, _ := repositoryProjector.NodeRepository.FindById(p2p.NodeId{Id:test.input.id})
		assert.Equal(t, node.GetID(), test.input.id)
		assert.Equal(t, node.IpAddress, test.input.address)
		repositoryProjector.NodeRepository.Remove(node.NodeId)
	}
}
func TestRepositoryProjector_HandleConnDisconnectedEvent(t *testing.T) {
	repositoryProjector, endUp := SetupRepositoryProjector("path_node_repository", "path_leader_repository")

	tests := map[string] struct{
		input struct{
			id string
		}
		err error
	}{
		"success":{
			input: struct {
				id      string
			}{id: string(123),},
			err: nil,
		},
		"empty node id test":{
			input: struct {
				id      string
			}{id: string(""),},
			err: adapter.ErrEmptyNodeId,
		},
	}

	defer endUp()

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		event := p2p.ConnectionDisconnectedEvent{
			EventModel: midgard.EventModel{
				ID:test.input.id,
			},
		}
		err := repositoryProjector.HandleConnDisconnectedEvent(event)
		assert.Equal(t, err, test.err)

		node, _ := repositoryProjector.NodeRepository.FindById(p2p.NodeId{Id:test.input.id})

		if node != nil{
			t.Errorf("node didn't removed!")
		}
	}
}
func TestRepositoryProjector_HandleLeaderUpdatedEvent(t *testing.T) {
	repositoryProjector, endUp := SetupRepositoryProjector("path_node_repository", "path_leader_repository")

	tests := map[string] struct{
		input struct{
			id string
		}
		err error
	}{
		"success":{
			input: struct {
				id      string
			}{id: string(123),},
			err: nil,
		},
		"empty node id test":{
			input: struct {
				id      string
			}{id: string(""),},
		},
		"empty address test":{
			input: struct {
				id      string
			}{id: string(123),},
		},
	}

	defer endUp()

	for testName, test := range tests{
		t.Logf("running test case %s", testName)
		event := p2p.LeaderUpdatedEvent{
			EventModel: midgard.EventModel{
				ID:test.input.id,
			},
		}
		err := repositoryProjector.HandleLeaderUpdatedEvent(event)
		assert.Equal(t, err, test.err)

		leader:= repositoryProjector.LeaderRepository.GetLeader()
		assert.Equal(t, leader.GetID(), test.input.id)
	}
}
func TestRepositoryProjector_HandlerNodeCreatedEvent(t *testing.T) {

}

func SetupRepositoryProjector(pathNodeRepository string, pathLeaderRepository string) (*adapter.RepositoryProjector, func()){
	nodeRepository := leveldb.NewNodeRepository(pathNodeRepository)
	leaderRepository := leveldb.NewLeaderRepository(pathLeaderRepository)

	repositoryProjector := adapter.NewRepositoryProjector(nodeRepository, leaderRepository)

	return repositoryProjector, func(){
		nodeRepository.Close()
		leaderRepository.Close()
	}
}