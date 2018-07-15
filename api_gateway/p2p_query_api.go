package api_gateway

import (
	"sync"

	"github.com/it-chain/it-chain-Engine/p2p"
)

type PLTableQueryApi struct {
	mux     sync.Mutex
	pLTable p2p.PLTable
}

func (pqa *PLTableQueryApi) FindById(peerId p2p.PeerId) (p2p.Peer, error) {

	pqa.mux.Lock()
	defer pqa.mux.Unlock()
	v, exist := pqa.peerTable[id.Id]

	if id.Id == "" {
		return v, ErrEmptyPeerId
	}
	//no matching id
	if !exist {
		return v, ErrNoMatchingPeer
	}

	return v, nil
}

func (pqa *PLTableQueryApi) FindAll() ([]p2p.Peer, error) {

}

func (pqa *PLTableQueryApi) FindByAddress(ipAddress string) (p2p.Peer, error) {

}

type PLTableQueryService interface {
	GetPLTable() (PLTable, error)
}
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
