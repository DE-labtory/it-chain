package p2p_test

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/magiconair/properties/assert"
)

func TestNewNode(t *testing.T) {
	node := p2p.NewNode("sdf",p2p.NodeId{Id:"123"})

	assert.Equal(t, node.IpAddress, "sdf")
	assert.Equal(t, node.NodeId.Id, "123")
}
