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
	blockApi BlockApi
}

func NewEventHandler(rp RepositoryProjector, api BlockApi) *EventHandler{
	return &EventHandler{
		repositoryProjector: rp,
		blockApi: api,
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

	syncState := blockchain.NewBlockSyncState()
	eh.repositoryProjector.EventRepository.Load(syncState, blockchain.BC_SYNC_STATE_AID)

	// TODO: sync state에 따라서 BlockApi 호출 여부 결정
	if !syncState.IsProgressing() {
		err := eh.blockApi.CheckAndSaveBlockFromPool(block.GetHeight())

		if err != nil {
			return err
		}
	}

	return nil
}
