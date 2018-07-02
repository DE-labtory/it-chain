package adapter

import (
	"github.com/it-chain/midgard"
	"github.com/it-chain/it-chain-Engine/blockchain"
	"errors"
)


var ErrEmptyEventId = errors.New("empty event id proposed.")
var ErrNodeApi = errors.New("problem in node api")


type RepositoryProjector struct {
	blockchain.PeerRepository
	midgard.EventRepository
}

type EventHandler struct {
	repositoryProjector RepositoryProjector
}

func NewEventHandler(rp RepositoryProjector) *EventHandler{
	return &EventHandler{
		repositoryProjector: rp,
	}
}

/// Check 단계에서 임의의 노드를 선정하기 위해 노드를 저장한다.
func (eh *EventHandler) HandleNodeCreatedEvent(event blockchain.NodeCreatedEvent) error {
	eventID := event.GetID()

	if eventID == "" {
		return ErrEmptyEventId
	}

	peer := event.Peer

	err := eh.repositoryProjector.PeerRepository.Add(peer)

	if err != nil {
		return ErrNodeApi
	}

	return nil
}

func (eh *EventHandler) HandleNodeDeletedEvent(event blockchain.NodeDeletedEvent) error {
	eventID := event.GetID()

	if eventID == "" {
		return ErrEmptyEventId
	}

	peer := event.Peer

	err := eh.repositoryProjector.PeerRepository.Remove(peer.PeerId)

	if err != nil {
		return ErrNodeApi
	}

	return nil
}

func (eh *EventHandler) HandleBlockQueuedEvent(event blockchain.BlockQueuedEvent) error {
	block := event.Block
	if block == nil {
		return ErrBlockNil
	}

	eh.repositoryProjector.EventRepository.Load()
	return nil
}
