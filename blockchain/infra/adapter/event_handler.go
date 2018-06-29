package adapter

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"errors"
)


var ErrEmptyEventId = errors.New("empty event id proposed.")
var ErrNodeApi = errors.New("problem in node api")


type RepositoryProjector struct {
	blockchain.PeerRepository
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

	node := event.Peer

	err := eh.repositoryProjector.PeerRepository.Add(node)

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

	node := event.Peer

	err := eh.repositoryProjector.PeerRepository.Remove(node.PeerId)

	if err != nil {
		return ErrNodeApi
	}

	return nil
}

