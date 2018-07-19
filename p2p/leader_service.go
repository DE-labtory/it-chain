package p2p

import (
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type ILeaderService interface {

	Set(leader Leader) error
}

type LeaderService struct {}

func (ls *LeaderService) Set(leader Leader) error{

	event := LeaderUpdatedEvent{
		EventModel:midgard.EventModel{
			ID:leader.LeaderId.Id,
		},
	}

	return eventstore.Save(leader.LeaderId.Id, event)
}
