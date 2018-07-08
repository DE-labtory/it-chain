package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/magiconair/properties/assert"
)


func TestCommunicationApi_DialToUnConnectedNode(t *testing.T) {
	tests := map[string]struct {
		input struct {
			peerList []p2p.Peer
		}
		err error
	}{
		"success": {
			input: struct{ peerList []p2p.Peer }{
				peerList: []p2p.Peer{
					{
						PeerId:p2p.PeerId{
							Id:"1",
						},
					},
				},
			},
			err:nil,
		},
	}

	for testName, test := range tests {

		t.Logf("running test case %s", testName)

		communicationApi := api.CommunicationApi{}

		assert.Equal(t, communicationApi.DialToUnConnectedNode(test.input.peerList), test.err)

	}

	mockPeerService := p2p.MockPeerService{}

	mockPeerService.FindByIdFunc = func(peerId p2p.PeerId) (p2p.Peer, error) {

		peerList := MakeFakePeerList()

		for _, peer := range peerList{
			if peer.PeerId == peerId{
				return peer, nil
			}
		}

		return p2p.Peer{}, nil
	}
}

func MakeFakePeerList() []p2p.Peer{

	peerList := make([]p2p.Peer, 0)
	peerList = append(peerList, p2p.Peer{
		PeerId:p2p.PeerId{
			Id:"1",
		},
		IpAddress:"1",
	})

	peerList = append(peerList, p2p.Peer{
		PeerId:p2p.PeerId{
			Id:"2",
		},
		IpAddress:"2",
	})

	return peerList
}