package adapter

import (
	"github.com/it-chain/engine/common/event"
)

type ConsensusCreatedEventHandler struct {
	//consensusApi api.ConsensusApi
}

func NewConsensusCreatedEventHandler() *ConsensusCreatedEventHandler {
	return &ConsensusCreatedEventHandler{}
}

func (handler ConsensusCreatedEventHandler) HandleConsensusCreatedEvent(e event.ConsensusCreated) {

}
