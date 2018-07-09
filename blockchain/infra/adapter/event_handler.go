package adapter

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/blockchain"
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

	err := eh.blockApi.CheckAndSaveBlockFromPool(height)

	if err != nil {
		return err
	}

	return nil
}
