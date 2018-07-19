package p2p

import (
	"errors"
	"fmt"

	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type Leader struct {
	LeaderId LeaderId
}

type LeaderId struct {
	Id string
}

func (lid LeaderId) ToString() string {
	return string(lid.Id)
}

func (l Leader) GetID() string {
	return l.LeaderId.ToString()
}

func (l *Leader) On(event midgard.Event) error {

	switch v := event.(type) {

	case LeaderChangedEvent:
		l.LeaderId = LeaderId{v.GetID()}

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

func UpdateLeader(peer Peer) error {

	leader := Leader{
		LeaderId: LeaderId{Id: peer.PeerId.Id},
	}

	if leader.LeaderId.Id == "" {
		return ErrEmptyLeaderId
	}

	events := make([]midgard.Event, 0)

	leaderUpdatedEvent := LeaderUpdatedEvent{
		EventModel: midgard.EventModel{
			ID:   leader.LeaderId.ToString(),
			Type: "leader.update",
		},
	}

	leader.On(leaderUpdatedEvent)

	events = append(events, leaderUpdatedEvent)

	return eventstore.Save(leaderUpdatedEvent.GetID(), events...)

}
