package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrEmptyPeerList = errors.New("empty peer list proposed")

type PeerApi struct {
	peerService          p2p.IPeerService
	peerQueryService     p2p.PeerQueryService
	pLTableQueryService  p2p.PLTableQueryService
	leaderService        p2p.LeaderService
	communicationService p2p.ICommunicationService
}

func NewPeerApi(
	peerService p2p.IPeerService,
	peerQueryService p2p.PeerQueryService,
	pLTableQueryService p2p.PLTableQueryService,
	leaderService p2p.LeaderService,
	communicationService p2p.ICommunicationService) *PeerApi {

	return &PeerApi{
		peerService:          peerService,
		peerQueryService:peerQueryService,
		pLTableQueryService:pLTableQueryService,
		leaderService:        leaderService,
		communicationService: communicationService,
	}
}

func (peerApi *PeerApi) UpdatePeerList(peerList []p2p.Peer) error {

	//둘다 존재할경우 무시, existPeerList에만 존재할경우 PeerDeletedEvent, peerList에 존재할경우 PeerCreatedEvent
	existPeerList, err := peerApi.peerQueryService.FindAll()

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

//Deliver Peer leader table that consists of peerList and leader
func (peerApi *PeerApi) DeliverPLTable(connectionId string) error {

	peerTable, _ := peerApi.pLTableQueryService.GetPLTable()
	peerApi.communicationService.DeliverPLTable(connectionId, peerTable)

	return nil
}

func (peerApi *PeerApi) FindById(id p2p.PeerId) (p2p.Peer, error) {

	return peerApi.peerQueryService.FindById(id)
}
