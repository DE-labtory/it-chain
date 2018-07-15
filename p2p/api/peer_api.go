package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrEmptyPeerList = errors.New("empty peer list proposed")

type PeerApi struct {
	peerQueryService     p2p.PeerQueryService
	pLTableQueryService  p2p.PLTableQueryService
	communicationService p2p.ICommunicationService
}

func NewPeerApi(
	peerQueryService p2p.PeerQueryService,
	pLTableQueryService p2p.PLTableQueryService,
	communicationService p2p.ICommunicationService) *PeerApi {

	return &PeerApi{
		peerQueryService:peerQueryService,
		pLTableQueryService:pLTableQueryService,
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

	//1. get peer table
	peerTable, _ := peerApi.pLTableQueryService.GetPLTable()

	//2. deliver peer table
	peerApi.communicationService.DeliverPLTable(connectionId, peerTable)

	return nil
}

func (peerApi *PeerApi) FindById(peerId p2p.PeerId) (p2p.Peer, error) {

	if peerId.Id==""{
		return p2p.Peer{PeerId:p2p.PeerId{Id:""}, IpAddress:""}, p2p.ErrEmptyPeerId
	}

	return peerApi.peerQueryService.FindById(peerId)
}
