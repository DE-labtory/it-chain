package memory

import (
	"errors"

	"sync"

	"github.com/it-chain/it-chain-Engine/p2p"
)

//types of errors
var ErrExistPeer = errors.New("proposed peer already exists")
var ErrNoMatchingPeer = errors.New("no matching peer exists")
var ErrEmptyPeerId = errors.New("empty peer id proposed")
var ErrEmptyPeerTable = errors.New("peer table is empty")

type PeerRepository struct {
	peerTable map[string]p2p.Peer
	mux       sync.Mutex
}

// 새로운 p2p repo 생성
func NewPeerRepository() (PeerRepository, error) {
	return PeerRepository{
		peerTable: make(map[string]p2p.Peer),
	}, nil
}


// done in peer service
func (pr *PeerRepository) Remove(id p2p.PeerId) error {

	pr.mux.Lock()
	defer pr.mux.Unlock()

	if id.Id == "" {
		return ErrEmptyPeerId
	}

	_, exist := pr.peerTable[id.Id]

	if !exist {
		return ErrNoMatchingPeer
	}

	delete(pr.peerTable, id.Id)

	return nil
}
