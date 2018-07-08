package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrEmptyPeerList = errors.New("empty peer list proposed")

type PeerApiCommunicationService interface{
	DeliverPLTable(connectionId string, peerLeaderTable p2p.PLTable) error
}

type PeerApi struct {
	peerService p2p.PeerService
	leaderService p2p.LeaderService
	communicationService PeerApiCommunicationService
}

func NewPeerApi(
	peerService p2p.PeerService,
	leaderService p2p.LeaderService,
	communicationService PeerApiCommunicationService) *PeerApi {

	return &PeerApi{
		peerService:     peerService,
		leaderService: leaderService,
		communicationService: communicationService,
	}
}

func (peerApi *PeerApi) UpdatePeerList(peerList []p2p.Peer) error {

	//둘다 존재할경우 무시, existPeerList에만 존재할경우 PeerDeletedEvent, peerList에 존재할경우 PeerCreatedEvent
	existPeerList, err := peerApi.peerService.FindAll()

	if err != nil {
		return err
	}

	newPeers, disconnectedPeers := p2p.GetMutuallyExclusivePeers(peerList, existPeerList)

	for _, peer := range newPeers {

		p2p.NewPeer(peer.IpAddress, peer.PeerId)
	}

	for _, peer := range disconnectedPeers {

		p2p.DeletePeer(peer.PeerId)
	}

	return nil
}

func (peerApi *PeerApi) GetPLTable() p2p.PLTable {

	leader := peerApi.leaderService.Get()
	peerList, _ := peerApi.peerService.FindAll()

	peerLeaderTable := p2p.PLTable{
		Leader:   leader,
		PeerList: peerList,
	}

	return peerLeaderTable
}

//Deliver Peer leader table that consists of peerList and leader
func (peerApi *PeerApi) DeliverPLTable(connectionId string) error {

	peerTable := peerApi.GetPLTable()
	peerApi.communicationService.DeliverPLTable(connectionId, peerTable)

	return nil
}

func (peerApi *PeerApi) FindById(id p2p.PeerId) (p2p.Peer, error) {

	return peerApi.peerService.FindById(id)
}
