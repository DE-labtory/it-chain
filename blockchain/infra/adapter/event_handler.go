package adapter

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
)

var ErrEmptyEventId = errors.New("empty event id proposed.")
var ErrNodeApi = errors.New("problem in node api")

type RepositoryProjector struct {
	PeerRepository blockchain.PeerRepository
}

/// Check 단계에서 임의의 노드를 선정하기 위해 노드를 저장한다.
func (eh *RepositoryProjector) HandleNodeCreatedEvent(event blockchain.NodeCreatedEvent) error {
	eventID := event.GetID()

	if eventID == "" {
		return ErrEmptyEventId
	}

	peer := event.Peer

	err := eh.PeerRepository.Add(peer)

	if err != nil {
		return ErrNodeApi
	}

	return nil
}

func (eh *RepositoryProjector) HandleNodeDeletedEvent(event blockchain.NodeDeletedEvent) error {
	eventID := event.GetID()

	if eventID == "" {
		return ErrEmptyEventId
	}

	peer := event.Peer

	err := eh.PeerRepository.Remove(peer.PeerId)

	if err != nil {
		return ErrNodeApi
	}

	return nil
}

type EventHandler struct {
	blockApi BlockApi
}

func NewEventHandler(api BlockApi) *EventHandler {
	return &EventHandler{
		blockApi: api,
	}
}

func (eh *EventHandler) HandleBlockAddToPoolEvent(event blockchain.BlockAddToPoolEvent) error {

	block := event.Block

	if block == nil {
		return ErrBlockNil
	}

	syncState := blockchain.NewBlockSyncState()
	eventstore.Load(syncState, blockchain.BC_SYNC_STATE_AID)

	// TODO: sync state에 따라서 BlockApi 호출 여부 결정
	if !syncState.IsProgressing() {
		err := eh.blockApi.CheckAndSaveBlockFromPool(block)

		if err != nil {
			return err
		}
	}

	return nil
}
