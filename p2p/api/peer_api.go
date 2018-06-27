package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
)

var ErrEmptyPeerList = errors.New("empty peer list proposed")

type ReadOnlyPeerRepository interface {
	FindById(id p2p.PeerId) (*p2p.Peer, error)
	FindAll() ([]p2p.Peer, error)
}

type PeerApi struct {
	peerRepository  ReadOnlyPeerRepository
	eventRepository EventRepository
	messageService  PeerMessageService
}

type PeerMessageService interface {
	DeliverPeerList(peerId p2p.PeerId, peerList []p2p.Peer) error
}

func NewPeerApi(peerRepository ReadOnlyPeerRepository, eventRepository EventRepository, messageService PeerMessageService) *PeerApi {
	return &PeerApi{
		peerRepository:  peerRepository,
		eventRepository: eventRepository,
		messageService:  messageService,
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

func (peerApi *PeerApi) DeliverPeerList(peerId p2p.PeerId) error {

	peerList, _ := peerApi.peerRepository.FindAll()
	if len(peerList) == 0 {
		return ErrEmptyPeerList
	}
	peerApi.messageService.DeliverPeerList(peerId, peerList)
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
