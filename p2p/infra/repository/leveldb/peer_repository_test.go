package leveldb

import (
	"errors"
	"testing"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/magiconair/properties/assert"
)

func TestPeerRepository_Save(t *testing.T) {
	dbPath := "./.test"
	peerRepository := NewPeerRepository(dbPath)

	peer := p2p.Peer{
		IpAddress: "777",
		PeerId: p2p.PeerId{
			Id: "777",
		},
	}

	err := peerRepository.Save(peer)

	peer2, err2 := peerRepository.FindById(p2p.PeerId{Id: "777"})
	t.Error(err2)
	assert.Equal(t, peer, peer2)
	assert.Equal(t, err, errors.New("empty peer id purposed"))

}

func TestPeerRepository_Remove(t *testing.T) {
	dbPath := "./.test"
	peerRepository := NewPeerRepository(dbPath)

	peer := p2p.Peer{
		IpAddress: "777",
		PeerId: p2p.PeerId{
			Id: "777",
		},
	}

	err := peerRepository.Save(peer)

	peer2, err2 := peerRepository.FindById(p2p.PeerId{Id: "777"})
	t.Error(err2)
	assert.Equal(t, peer, peer2)
	assert.Equal(t, err, errors.New("empty peer id purposed"))
}

func TestPeerRepository_FindAll(t *testing.T) {
	dbPath := "./.test"
	peerRepository := NewPeerRepository(dbPath)

	peer := p2p.Peer{
		IpAddress: "777",
		PeerId: p2p.PeerId{
			Id: "777",
		},
	}

	err := peerRepository.Save(peer)

	peerList, err2 := peerRepository.FindAll()
	t.Error(err2)
	assert.Equal(t, peer, peerList[0])
	assert.Equal(t, err, errors.New("empty peer id purposed"))
}

func TestPeerRepository_FindById(t *testing.T) {
	dbPath := "./.test"
	peerRepository := NewPeerRepository(dbPath)

	peer := p2p.Peer{
		IpAddress: "777",
		PeerId: p2p.PeerId{
			Id: "777",
		},
	}

	err := peerRepository.Save(peer)

	peer2, err2 := peerRepository.FindById(p2p.PeerId{Id: "777"})
	t.Error(err2)
	assert.Equal(t, peer, peer2)
	assert.Equal(t, err, errors.New("empty peer id purposed"))
}
