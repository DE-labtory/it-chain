package api_gateway

import (
	"sync"

	"github.com/it-chain/it-chain-Engine/p2p"
)

type PLTableQueryApi struct {
	mux     sync.Mutex
	pLTable p2p.PLTable
}

func (pltqa *PLTableQueryApi) GetPLTable() (p2p.PLTable, error){

	return pltqa.pLTable, nil
}

func (pltqa *PLTableQueryApi) FindPeerById(peerId p2p.PeerId) (p2p.Peer, error) {

	pltqa.mux.Lock()
	defer pltqa.mux.Unlock()
	v, exist := pltqa.pLTable.PeerTable[peerId.Id]

	if peerId.Id == "" {
		return v, p2p.ErrEmptyPeerId
	}
	//no matching id
	if !exist {
		return v, p2p.ErrNoMatchingPeerId
	}

	return v, nil
}

func (pltqa *PLTableQueryApi) FindPeerByAddress(ipAddress string) (p2p.Peer, error) {

	pltqa.mux.Lock()
	defer pltqa.mux.Unlock()

	for _, peer := range pltqa.pLTable.PeerTable {

		if peer.IpAddress == ipAddress {
			return peer, nil
		}
	}

	return p2p.Peer{}, nil
}

type PeerRepository struct {
	peerTable map[string]p2p.Peer
	mux       sync.Mutex
}

