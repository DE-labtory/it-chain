package api_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
)

type MockCommunicationServiceForDial struct {
	dialFunc func(ipAddress string) error
}

func TestCommunicationApi_DialToUnConnectedNode(t *testing.T) {
	tests := map[string]struct {
		input struct {
			peerList []p2p.Peer
		}
	}{
		"empty peer list test": {
			input: struct{ peerList []p2p.Peer }{peerList: nil},
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

	}

	mockPeerService := p2p.MockPeerService{}

	mockPeerService.FindByIdFunc = func(peerId p2p.PeerId) (p2p.Peer, error) {
		peer := p2p.Peer{
			PeerId:peerId,
		}
	}
}
