package adapter

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
)

var ErrEmptyEventId = errors.New("empty event id proposed.")
var ErrNodeApi = errors.New("problem in node api")

type EventHandler struct {
	blockApi BlockApi
}

func NewEventHandler(api BlockApi) *EventHandler {
	return &EventHandler{
		blockApi: api,
	}
}

func (eh *EventHandler) HandleBlockAddToPoolEvent(event blockchain.BlockAddToPoolEvent) error {
	height := event.Height

	if height < 0 {
		return ErrBlockNil
	}

	syncState := blockchain.NewBlockSyncState()
	eventstore.Load(syncState, blockchain.BC_SYNC_STATE_AID)

	// TODO: sync state에 따라서 BlockApi 호출 여부 결정
	if !syncState.IsProgressing() {
		err := eh.blockApi.CheckAndSaveBlockFromPool(height)

		if err != nil {
			return err
		}
	}

	return nil
}
