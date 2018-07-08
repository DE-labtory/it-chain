package p2p

import (
	"errors"
	"fmt"

	"github.com/it-chain/midgard"
	"log"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
)

type Leader struct {
	LeaderId LeaderId
}

type LeaderId struct {
	Id string
}

type LeaderService interface {

	Get() Leader
	Set(leader Leader) error
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

func UpdateLeader(leader Leader) error{

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

	err := eventstore.Save(leaderUpdatedEvent.GetID(), events...)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return err
}

type LeaderRepository interface {
	GetLeader() Leader
	SetLeader(leader Leader)
}
