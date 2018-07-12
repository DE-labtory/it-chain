package api_test

import (
	"errors"
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/magiconair/properties/assert"
)

var ErrEmptyPeerList = errors.New("empty peer list proposed")

//todo make node api test
//todo make fake dependencies 1. eventRepository 2. messageDispatcher 3. peerRepository
//todo make test map
//todo test continue

type MockService struct{}

func (ms MockService) GetPeerLeaderTable() p2p.PeerLeaderTable {
	peerLeaderTable := p2p.PeerLeaderTable{
		Leader: p2p.Leader{
			LeaderId: p2p.LeaderId{Id: "1"},
		},
		PeerList: []p2p.Peer{{
			PeerId: p2p.PeerId{
				Id: "2",
			},
		}},
	}
	return peerLeaderTable
}

type MockPeerRepository struct{}

func (mnr MockPeerRepository) FindById(id p2p.PeerId) (p2p.Peer, error) {
	peer := p2p.Peer{PeerId: id}
	return peer, nil
}
func (mnr MockPeerRepository) FindAll() ([]p2p.Peer, error) { return nil, nil }

type MockLeaderRepository struct{}

func (mpr MockLeaderRepository) GetLeader() p2p.Leader {
	leader := p2p.Leader{LeaderId: p2p.LeaderId{Id: "1"}}
	return leader
}

type MockPeerApiGrpcCommandService struct{}

func (mnms MockPeerApiGrpcCommandService) DeliverPeerLeaderTable(connectionId string, peerLeaderTable p2p.PeerLeaderTable) error {
	return nil
}

func TestPeerApi_UpdatePeerList(t *testing.T) {

	tests := map[string]struct {
		input []p2p.Peer
		err   error
	}{
		"success": {
			input: []p2p.Peer{
				p2p.Peer{
					PeerId: p2p.PeerId{
						Id: "1",
					},
				},
			},
			err: nil,
		},
	}

	peerApi := SetupPeerApi()

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := peerApi.UpdatePeerList(test.input)
		assert.Equal(t, err, test.err)
	}
}

func SetupPeerApi() *api.PeerApi {
	mockService := MockService{}
	peerRepository := MockPeerRepository{}
	leaderRepository := MockLeaderRepository{}
	eventRepository := MockEventRepository{}
	grpcCommandService := MockPeerApiGrpcCommandService{}

	peerApi := api.NewPeerApi(mockService, peerRepository, leaderRepository, eventRepository, grpcCommandService)

	return peerApi
}
