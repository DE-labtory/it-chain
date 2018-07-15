package api_gateway

import (
	"sync"

	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrPeerExists = errors.New("peer already exists")

type PLTableQueryApi struct {
	mux               sync.Mutex
	pLTableRepository PLTableRepository
}

func (pltqa *PLTableQueryApi) GetPLTable() (p2p.PLTable, error) {

	return pltqa.pLTableRepository.GetPLTable()
}

func (pltqa *PLTableQueryApi) GetLeader() (p2p.Leader, error) {

	return pltqa.pLTableRepository.GetLeader()
}

func (pltqa *PLTableQueryApi) FindPeerById(peerId p2p.PeerId) (p2p.Peer, error) {

	return pltqa.pLTableRepository.FindPeerById(peerId)
}

func (pltqa *PLTableQueryApi) FindPeerByAddress(ipAddress string) (p2p.Peer, error) {

	return pltqa.pLTableRepository.FindPeerByAddress(ipAddress)
}

type PLTableRepository struct {
	mux     sync.Mutex
	pLTable p2p.PLTable
}

func (pltrepo *PLTableRepository) GetPLTable() (p2p.PLTable, error) {

	pltrepo.mux.Lock()
	defer pltrepo.mux.Unlock()

	return pltrepo.pLTable, nil
}

func (pltrepo *PLTableRepository) GetLeader() (p2p.Leader, error) {

	return pltrepo.pLTable.Leader, nil
}

func (pltrepo *PLTableRepository) FindPeerById(peerId p2p.PeerId) (p2p.Peer, error) {

	pltrepo.mux.Lock()
	defer pltrepo.mux.Unlock()
	v, exist := pltrepo.pLTable.PeerTable[peerId.Id]

	if peerId.Id == "" {
		return v, p2p.ErrEmptyPeerId
	}
	//no matching id
	if !exist {
		return v, p2p.ErrNoMatchingPeerId
	}

	return v, nil
}

func (pltrepo *PLTableRepository) FindPeerByAddress(ipAddress string) (p2p.Peer, error) {

	pltrepo.mux.Lock()
	defer pltrepo.mux.Unlock()

	for _, peer := range pltrepo.pLTable.PeerTable {

		if peer.IpAddress == ipAddress {
			return peer, nil
		}
	}

	return p2p.Peer{}, nil
}

func (pltrepo *PLTableRepository) Save(peer p2p.Peer) error {

	pltrepo.mux.Lock()
	defer pltrepo.mux.Unlock()

	pLTable, _ := pltrepo.GetPLTable()

	_, exist := pLTable.PeerTable[peer.PeerId.Id]

	if exist {
		return ErrPeerExists
	}

	pLTable.PeerTable[peer.PeerId.Id] = peer

	return nil
}

func (pltrepo *PLTableRepository) SetLeader(peer p2p.Peer) error{

	pltrepo.mux.Lock()
	defer pltrepo.mux.Unlock()

	leader := p2p.Leader{
		LeaderId:p2p.LeaderId{
			Id:peer.PeerId.Id,
		},
	}

	pltrepo.pLTable.Leader = leader

	return nil
}

type P2PEventHandler struct {
	pLTableRepository PLTableRepository
}

func (peh *P2PEventHandler) PeerCreatedEventHandler(event p2p.PeerCreatedEvent) error{

	peer := p2p.Peer{
		PeerId:p2p.PeerId{
			Id:event.ID,
		},
		IpAddress:event.IpAddress,
	}

	peh.pLTableRepository.Save(peer)

	return nil
}

func (peh *P2PEventHandler) HandleLeaderUpdatedEvent(event p2p.LeaderUpdatedEvent) error {

	peer := p2p.Peer{
		PeerId: p2p.PeerId{Id: event.ID},
	}

	peh.pLTableRepository.SetLeader(peer)

	return nil

}