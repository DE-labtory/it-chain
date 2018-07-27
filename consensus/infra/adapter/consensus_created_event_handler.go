package adapter

import (
	"github.com/it-chain/engine/consensus/api"
	"github.com/it-chain/engine/common/event"
)

type ConsensusCreatedEventHandler struct {
	consensusApi api.ConsensusApi
}

func NewConsensusCreatedEventHandler(consensusApi api.ConsensusApi) *ConsensusCreatedEventHandler{
	return &ConsensusCreatedEventHandler{
		consensusApi:consensusApi,
	}
}

func (handler ConsensusCreatedEventHandler) HandleConsensusCreatedEvent(e event.ConsensusCreated){

}
