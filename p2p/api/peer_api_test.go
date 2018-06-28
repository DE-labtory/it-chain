package api_test

import (
	"errors"
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/magiconair/properties/assert"
)

var ErrEmptyPeerList = errors.New("empty node list proposed")

//todo make node api test
//todo make fake dependencies 1. eventRepository 2. messageDispatcher 3. nodeRepository
//todo make test map
//todo test continue

type MockPeerRepository struct {}

func (mnr MockPeerRepository) FindById(id p2p.PeerId) (*p2p.Peer, error) { return nil, nil }
func (mnr MockPeerRepository) FindAll() ([]p2p.Peer, error)              { return nil, nil }

type MockPeerMessageService struct{}

func (mnms MockPeerMessageService) DeliverPeerList(nodeId p2p.PeerId, nodeList []p2p.Peer) error {
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

	nodeApi := SetupPeerApi()

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := nodeApi.UpdatePeerList(test.input)
		assert.Equal(t, err, test.err)
	}
}

func SetupPeerApi() *api.PeerApi {
	nodeRepository := MockPeerRepository{}
	eventRepository := MockEventRepository{}
	messageService := MockPeerMessageService{}

	nodeApi := api.NewPeerApi(nodeRepository, eventRepository, messageService)

	return nodeApi
}
