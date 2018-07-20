package adapter

import (
	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/engine/consensus/api"
)

type LeaderChangedEventHandler struct {
	parliamentApi api.ParliamentApi
}

func (handler LeaderChangedEventHandler) HandleLeaderChangedEvent(event consensus.LeaderChangedEvent) error {

	leaderId := consensus.LeaderId{event.LeaderId}

	err := handler.parliamentApi.ChangeLeader(consensus.Leader{leaderId})

	if err != nil {
		return err
	}
	return nil

}
