package api

import (
	"errors"
	"log"
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
)

var ErrEmptyPeerId = errors.New("empty peer id requested")
var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyConnectionId = errors.New("empty connection id proposed")

type LeaderApi struct {
	leaderRepository   ReadOnlyLeaderRepository
	eventRepository    EventRepository
	myInfo             *p2p.Peer
}

type Publish func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

type ReadOnlyLeaderRepository interface {
	GetLeader() p2p.Leader
}

type EventRepository interface { //midgard.Repository
	Save(aggregateID string, events ...midgard.Event) error
}

func NewLeaderApi(leaderRepository ReadOnlyLeaderRepository, eventRepository EventRepository, myInfo *p2p.Peer) *LeaderApi {

	return &LeaderApi{
		leaderRepository:   leaderRepository,
		eventRepository:    eventRepository,
		myInfo:             myInfo,
	}
}

//todo update leader with ip address by peer!! not leader
func (leaderApi *LeaderApi) UpdateLeader(leader p2p.Leader) error {

	if leader.LeaderId.Id == "" {
		return ErrEmptyLeaderId
	}

	events := make([]midgard.Event, 0)
	leaderUpdatedEvent := p2p.LeaderUpdatedEvent{
		EventModel: midgard.EventModel{
			ID:   leader.LeaderId.ToString(),
			Type: "leader.update",
		},
	}

	events = append(events, leaderUpdatedEvent)
	err := leaderApi.eventRepository.Save(leaderUpdatedEvent.GetID(), events...)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	return err
}
