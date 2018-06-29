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

type MockService struct {}
func (ms MockService) GetPeerTable() p2p.PeerTable{
	peerTable := p2p.PeerTable{
		Leader:p2p.Leader{
			LeaderId:p2p.LeaderId{Id:"1"},
		},
		PeerList:[]p2p.Peer{{
			PeerId:p2p.PeerId{
				Id:"2",
			},
		}},
	}
	return peerTable
}
type MockPeerRepository struct {}

func (mnr MockPeerRepository) FindById(id p2p.PeerId) (p2p.Peer, error) {
	peer := p2p.Peer{PeerId:id}
	return peer, nil
}
func (mnr MockPeerRepository) FindAll() ([]p2p.Peer, error)              { return nil, nil }

type MockPeerApiGrpcCommandService struct{}

func (mnms MockPeerApiGrpcCommandService) DeliverPeerTable(connectionId string, peerTable p2p.PeerTable) error {
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
	eventRepository := MockEventRepository{}
	grpcCommandService := MockPeerApiGrpcCommandService{}

	peerApi := api.NewPeerApi(mockService, peerRepository, eventRepository, grpcCommandService)

	return peerApi
}
