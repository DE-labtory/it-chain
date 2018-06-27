package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
)

var ErrEmptyPeerList = errors.New("empty node list proposed")

type ReadOnlyPeerRepository interface {
	FindById(id p2p.PeerId) (*p2p.Peer, error)
	FindAll() ([]p2p.Peer, error)
}

type PeerApi struct {
	nodeRepository  ReadOnlyPeerRepository
	eventRepository EventRepository
	messageService  PeerMessageService
}

type PeerMessageService interface {
	DeliverPeerList(nodeId p2p.PeerId, nodeList []p2p.Peer) error
}

func NewPeerApi(nodeRepository ReadOnlyPeerRepository, eventRepository EventRepository, messageService PeerMessageService) *PeerApi {
	return &PeerApi{
		nodeRepository:  nodeRepository,
		eventRepository: eventRepository,
		messageService:  messageService,
	}
}

func (nodeApi *PeerApi) UpdatePeerList(nodeList []p2p.Peer) error {

	//둘다 존재할경우 무시, existPeerList에만 존재할경우 PeerDeletedEvent, nodeList에 존재할경우 PeerCreatedEvent
	var event midgard.Event

	existPeerList, err := nodeApi.nodeRepository.FindAll()

	if err != nil {
		return err
	}

	newPeers, disconnectedPeers := p2p.GetMutuallyExclusivePeers(nodeList, existPeerList)

	for _, node := range newPeers {

		event = p2p.PeerCreatedEvent{
			EventModel: midgard.EventModel{
				ID:   node.GetID(),
				Type: "node.created",
			},
			IpAddress: node.IpAddress,
		}

		nodeApi.eventRepository.Save(event.GetID(), event)
	}

	for _, node := range disconnectedPeers {
		event = p2p.PeerDeletedEvent{
			EventModel: midgard.EventModel{
				ID:   node.GetID(),
				Type: "node.deleted",
			},
		}

		nodeApi.eventRepository.Save(event.GetID(), event)
	}

	return nil
}

func (nodeApi *PeerApi) DeliverPeerList(nodeId p2p.PeerId) error {

	nodeList, _ := nodeApi.nodeRepository.FindAll()
	if len(nodeList) == 0 {
		return ErrEmptyPeerList
	}
	nodeApi.messageService.DeliverPeerList(nodeId, nodeList)
	return nil
}

// add a node
func (nodeApi *PeerApi) AddPeer(node p2p.Peer) error {

	event := p2p.PeerCreatedEvent{
		EventModel: midgard.EventModel{
			ID:   node.GetID(),
			Type: "node.created",
		},
		IpAddress: node.IpAddress,
	}

	err := nodeApi.eventRepository.Save(event.GetID(), event)

	if err != nil {
		return err
	}

	return nil
}

// delete a node
func (nodeApi *PeerApi) DeletePeer(id p2p.PeerId) error {

	event := p2p.PeerDeletedEvent{
		EventModel: midgard.EventModel{
			ID:   id.ToString(),
			Type: "node.deleted",
		},
	}

	err := nodeApi.eventRepository.Save(event.GetID(), event)

	if err != nil {
		return err
	}

	return nil
}
