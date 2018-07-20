package adapter

import (
	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/engine/consensus/api"
)

type PrepareMsgAddedEventHandler struct {
	consensusApi api.ConsensusApi
}

func (handler PrepareMsgAddedEventHandler) HandlePrepareMsgAddedEvent(event consensus.PrepareMsgAddedEvent) error {
	err := handler.consensusApi.ReceivePrepareMsg(event.PrepareMsg)

	if err != nil {
		return err
	}
	return nil
}
