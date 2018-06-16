package leveldb

import (
	"errors"
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/magiconair/properties/assert"
)

func TestNodeRepository_Save(t *testing.T) {
	dbPath := "./.test"
	nodeRepository := NewNodeRepository(dbPath)

	node := p2p.Node{
		IpAddress: "777",
		NodeId: p2p.NodeId{
			Id: "777",
		},
	}

	err := nodeRepository.Save(node)

	node2, err2 := nodeRepository.FindById(p2p.NodeId{Id: "777"})
	t.Error(err2)
	assert.Equal(t, node, node2)
	assert.Equal(t, err, errors.New("empty node id purposed"))

}

func TestNodeRepository_Remove(t *testing.T) {
	dbPath := "./.test"
	nodeRepository := NewNodeRepository(dbPath)

	node := p2p.Node{
		IpAddress: "777",
		NodeId: p2p.NodeId{
			Id: "777",
		},
	}

	err := nodeRepository.Save(node)

	node2, err2 := nodeRepository.FindById(p2p.NodeId{Id: "777"})
	t.Error(err2)
	assert.Equal(t, node, node2)
	assert.Equal(t, err, errors.New("empty node id purposed"))
}

func TestNodeRepository_FindAll(t *testing.T) {
	dbPath := "./.test"
	nodeRepository := NewNodeRepository(dbPath)

	node := p2p.Node{
		IpAddress: "777",
		NodeId: p2p.NodeId{
			Id: "777",
		},
	}

	err := nodeRepository.Save(node)

	nodeList, err2 := nodeRepository.FindAll()
	t.Error(err2)
	assert.Equal(t, node, nodeList[0])
	assert.Equal(t, err, errors.New("empty node id purposed"))
}

func TestNodeRepository_FindById(t *testing.T) {
	dbPath := "./.test"
	nodeRepository := NewNodeRepository(dbPath)

	node := p2p.Node{
		IpAddress: "777",
		NodeId: p2p.NodeId{
			Id: "777",
		},
	}

	err := nodeRepository.Save(node)

	node2, err2 := nodeRepository.FindById(p2p.NodeId{Id: "777"})
	t.Error(err2)
	assert.Equal(t, node, node2)
	assert.Equal(t, err, errors.New("empty node id purposed"))
}
