package adapter

import (
	"github.com/it-chain/engine/common/event"
)

type PrepareMessageAddedEventHandler struct {
}

func NewPrepareMessageEventHandler() *PrepareMessageAddedEventHandler {
	return &PrepareMessageAddedEventHandler{}
}

func (handler PrepareMessageAddedEventHandler) HandlePrepareMessage(e event.PrepareMsgAdded) {

	// prepare message 가 들어왔을때 처리?
	//handler.consensusApi.ReceivePrepareMsg()
}
