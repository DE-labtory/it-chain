package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
)

var ErrEmptyPeerList = errors.New("empty peer list proposed")

type PeerApiGrpcCommandService interface {
	DeliverPLTable(connectionId string, peerTable p2p.PLTable) error
}
type ReadOnlyPeerRepository interface {
	FindById(id p2p.PeerId) (p2p.Peer, error)
	FindAll() ([]p2p.Peer, error)
}

type PeerApi struct {
	peerRepository            ReadOnlyPeerRepository
	leaderRepository          ReadOnlyLeaderRepository
	eventRepository           EventRepository
	peerApiGrpcCommandService PeerApiGrpcCommandService
}

func NewPeerApi(peerRepository ReadOnlyPeerRepository, leaderRepository ReadOnlyLeaderRepository, eventRepository EventRepository, peerApiGrpcCommandService PeerApiGrpcCommandService) *PeerApi {
	return &PeerApi{
		peerRepository:            peerRepository,
		leaderRepository:          leaderRepository,
		eventRepository:           eventRepository,
		peerApiGrpcCommandService: peerApiGrpcCommandService,
	}
}

func (peerApi *PeerApi) UpdatePeerList(peerList []p2p.Peer) error {

	//둘다 존재할경우 무시, existPeerList에만 존재할경우 PeerDeletedEvent, peerList에 존재할경우 PeerCreatedEvent
	var event midgard.Event

	existPeerList, err := peerApi.peerRepository.FindAll()

	if err != nil {
		return err
	}

	newPeers, disconnectedPeers := p2p.GetMutuallyExclusivePeers(peerList, existPeerList)

	for _, peer := range newPeers {

		event = p2p.PeerCreatedEvent{
			EventModel: midgard.EventModel{
				ID:   peer.GetID(),
				Type: "peer.created",
			},
			IpAddress: peer.IpAddress,
		}

		peerApi.eventRepository.Save(event.GetID(), event)
	}

	for _, peer := range disconnectedPeers {
		event = p2p.PeerDeletedEvent{
			EventModel: midgard.EventModel{
				ID:   peer.GetID(),
				Type: "peer.deleted",
			},
		}

		peerApi.eventRepository.Save(event.GetID(), event)
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

//Deliver Peer table that consists of peerList and leader
func (peerApi *PeerApi) DeliverPLTable(connectionId string) error {

	peerTable := peerApi.GetPLTable()
	peerApi.peerApiGrpcCommandService.DeliverPLTable(connectionId, peerTable)

	return nil
}

// add a peer
func (peerApi *PeerApi) AddPeer(peer p2p.Peer) (p2p.Peer, error) {

	if peer.PeerId.Id == ""{
		
		return p2p.Peer{}, ErrEmptyPeerId
	}

	return p2p.NewPeer(peer.IpAddress, peer.PeerId)
}

// delete a peer
func (peerApi *PeerApi) DeletePeer(id p2p.PeerId) error {

	if id.Id == ""{
		return ErrEmptyPeerId
	}

	return p2p.DeletePeer(id)
}

func (peerApi *PeerApi) FindById(id p2p.PeerId) (p2p.Peer, error) {

	return peerApi.peerRepository.FindById(id)
}
