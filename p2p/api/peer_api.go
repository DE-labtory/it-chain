package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
)

var ErrEmptyPeerList = errors.New("empty peer list proposed")
type PeerApiService interface {
	GetPeerTable() p2p.PeerTable
}
type PeerApiGrpcCommandService interface {
	DeliverPeerTable(connectionId string, peerTable p2p.PeerTable) error
}
type ReadOnlyPeerRepository interface {
	FindById(id p2p.PeerId) (p2p.Peer, error)
	FindAll() ([]p2p.Peer, error)
}

type PeerApi struct {
	service PeerApiService
	peerRepository  ReadOnlyPeerRepository
	eventRepository EventRepository
	peerApiGrpcCommandService  PeerApiGrpcCommandService
}



func NewPeerApi(service PeerApiService, peerRepository ReadOnlyPeerRepository, eventRepository EventRepository, peerApiGrpcCommandService PeerApiGrpcCommandService) *PeerApi {
	return &PeerApi{
		service:service,
		peerRepository:  peerRepository,
		eventRepository: eventRepository,
		peerApiGrpcCommandService:  peerApiGrpcCommandService,
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
func (peerApi *PeerApi) GetPeerTable() p2p.PeerTable{
	peerTable := peerApi.service.GetPeerTable()

	return peerTable
}

//Deliver Peer table that consists of peerList and leader
func (peerApi *PeerApi) DeliverPeerTable(connectionId string) error {

	peerTable := peerApi.service.GetPeerTable()
	peerApi.peerApiGrpcCommandService.DeliverPeerTable(connectionId, peerTable)
	return nil
}

// add a peer
func (peerApi *PeerApi) AddPeer(peer p2p.Peer) error {

	event := p2p.PeerCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   peer.GetID(),
			Type: "peer.created",
		},
		IpAddress: peer.IpAddress,
	}

	err := peerApi.eventRepository.Save(event.GetID(), event)

	if err != nil {
		return err
	}

	return nil
}

// delete a peer
func (peerApi *PeerApi) DeletePeer(id p2p.PeerId) error {

	event := p2p.PeerDeletedEvent{
		EventModel: midgard.EventModel{
			ID:   id.ToString(),
			Type: "peer.deleted",
		},
	}

	err := peerApi.eventRepository.Save(event.GetID(), event)

	if err != nil {
		return err
	}

	return nil
}
