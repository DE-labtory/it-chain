package adapter

import (
	"github.com/it-chain/engine/consensus/api"
	"github.com/it-chain/engine/consensus"
)

type CommitMsgEventHandler struct {
	consensusApi api.ConsensusApi
}

func (handler CommitMsgEventHandler) HandleCommitMsgEvent(event consensus.CommitMsgAddedEvent) error {

	err := handler.consensusApi.ReceiveCommitMsg(event.CommitMsg)

	if err != nil {
		return err
	}
	return nil

}
