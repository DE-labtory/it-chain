package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrEmptyPeerList = errors.New("empty peer list proposed")

type ReadOnlyPeerRepository interface {
	FindById(id p2p.PeerId) (p2p.Peer, error)
	FindAll() ([]p2p.Peer, error)
}

type PeerApiCommunicationService interface{
	DeliverPLTable(connectionId string, peerLeaderTable p2p.PLTable) error
}

type PeerApi struct {
	peerRepository       ReadOnlyPeerRepository
	leaderRepository     ReadOnlyLeaderRepository
	communicationService PeerApiCommunicationService
}

func NewPeerApi(
	peerRepository ReadOnlyPeerRepository,
	leaderRepository ReadOnlyLeaderRepository,
	communicationService PeerApiCommunicationService) *PeerApi {

	return &PeerApi{
		peerRepository:       peerRepository,
		leaderRepository:     leaderRepository,
		communicationService: communicationService,
	}
}

func (peerApi *PeerApi) UpdatePeerList(peerList []p2p.Peer) error {

	//둘다 존재할경우 무시, existPeerList에만 존재할경우 PeerDeletedEvent, peerList에 존재할경우 PeerCreatedEvent
	existPeerList, err := peerApi.peerRepository.FindAll()

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

	leader := peerApi.leaderRepository.GetLeader()
	peerList, _ := peerApi.peerRepository.FindAll()

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

	return peerApi.peerRepository.FindById(id)
}
