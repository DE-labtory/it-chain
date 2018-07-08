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

// done in peer service named save
func (pr *PeerRepository) Save(data p2p.Peer) error {

	pr.mux.Lock()
	defer pr.mux.Unlock()

	// return empty peerID error if peerID is null
	if data.PeerId.Id == "" {
		return ErrEmptyPeerId
	}
	_, exist := pr.peerTable[data.PeerId.Id]

	if exist {
		return ErrExistPeer
	}

	pr.peerTable[data.PeerId.Id] = data

	return nil
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

// done in peer service
func (pr *PeerRepository) FindById(id p2p.PeerId) (p2p.Peer, error) {
	pr.mux.Lock()
	defer pr.mux.Unlock()
	v, exist := pr.peerTable[id.Id]

	if id.Id == "" {
		return v, ErrEmptyPeerId
	}
	//no matching id
	if !exist {
		return v, ErrNoMatchingPeer
	}

	return v, nil
}

// done in peer service
func (pr *PeerRepository) FindAll() ([]p2p.Peer, error) {
	pr.mux.Lock()
	defer pr.mux.Unlock()
	peers := make([]p2p.Peer, 0)

	if len(pr.peerTable) == 0 {
		return peers, ErrEmptyPeerTable
	}

	for _, value := range pr.peerTable {
		peers = append(peers, value)
	}

	return peers, nil
}


func (pr *PeerRepository) FindByAddress(ipAddress string) (p2p.Peer, error) {

	pr.mux.Lock()
	defer pr.mux.Unlock()

	for _, peer := range pr.peerTable {

		if peer.IpAddress == ipAddress {
			return peer, nil
		}
	}

	return p2p.Peer{}, nil
}
