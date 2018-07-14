package adapter

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
)

type EventHandler struct {
	blockApi BlockApi
}

func NewEventHandler(api BlockApi) *EventHandler {
	return &EventHandler{
		blockApi: api,
	}
}

// TODO: write test case
func (eh *EventHandler) HandleBlockAddToPoolEvent(event blockchain.BlockStagedEvent) error {
	blockId := event.ID

	if blockId == "" {
		return ErrBlockNil
	}

	err := eh.blockApi.CommitBlockFromPoolOrSync(blockId)

	if err != nil {
		return err
	}

	return nil
}
