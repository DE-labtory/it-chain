package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
	"github.com/magiconair/properties/assert"
	"github.com/it-chain/it-chain-Engine/p2p/test/mock"
)

func TestCommunicationApi_DialToUnConnectedNode(t *testing.T) {
	tests := map[string]struct {
		input struct {
			peerTable map[string]p2p.Peer
		}
		err error
	}{
		"success": {
			input: struct{ peerTable map[string]p2p.Peer }{peerTable: map[string]p2p.Peer{
				"1":{
					PeerId:p2p.PeerId{
						Id:"1",
					},
				},
			}},
			err: nil,
		},
	}

	mockPLTableQueryService :=&mock.MockPLTableQueryService{}
	mockPLTableQueryService.FindPeerByIdFunc = func(peerId p2p.PeerId) (p2p.Peer, error) {

		peerTable := mock.MakeFakePeerTable()

		for _, peer := range peerTable {
			if peer.PeerId == peerId {
				return peer, nil
			}
		}

		return p2p.Peer{}, nil
	}

	for testName, test := range tests {

		t.Logf("running test case %s", testName)

		communicationApi := api.NewCommunicationApi(mockPLTableQueryService, &mock.MockCommunicationService{})

		assert.Equal(t, communicationApi.DialToUnConnectedNode(test.input.peerTable), test.err)

	}
}
