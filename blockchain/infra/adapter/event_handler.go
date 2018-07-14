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
func (eh *EventHandler) HandleBlockAddToPoolEvent(event blockchain.BlockAddToPoolEvent) error {
	height := event.Height

	if height < 0 {
		return ErrBlockNil
	}

	err := eh.blockApi.CheckAndSaveBlockFromPool(height)

	if err != nil {
		return err
	}

	return nil
}
