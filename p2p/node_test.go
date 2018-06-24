package p2p_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
)

func TestNewNode(t *testing.T) {

	node := p2p.NewNode("sdf", p2p.NodeId{Id: "sdf"})

	if node.IpAddress != "sdf" {
		t.Error("new node failed!")
	}
	if node.NodeId.Id != "sdf" {
		t.Error("new node failed!")
	}
}
