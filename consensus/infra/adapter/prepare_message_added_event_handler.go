package adapter

import (
	"github.com/it-chain/engine/consensus/api"
	"github.com/it-chain/engine/common/event"
)

type PrepareMessageAddedEventHandler struct {
	consensusApi api.ConsensusApi
}

func NewPrepareMessageEventHandler(consensusApi api.ConsensusApi) *PrepareMessageAddedEventHandler{

	return &PrepareMessageAddedEventHandler{
		consensusApi: consensusApi,
	}
}

func (handler PrepareMessageAddedEventHandler) HandlePrepareMessage(e event.PrepareMsgAdded){

	// prepare message 가 들어왔을때 처리?
	handler.consensusApi.ReceivePrepareMsg()
}
