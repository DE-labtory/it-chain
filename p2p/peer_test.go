package p2p_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
)

func TestNewPeer(t *testing.T) {

	peer := p2p.NewPeer("sdf", p2p.PeerId{Id: "sdf"})

	if peer.IpAddress != "sdf" {
		t.Error("new peer failed!")
	}
	if peer.PeerId.Id != "sdf" {
		t.Error("new peer failed!")
	}
}
