package api

import (
	"log"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
	"errors"
)

type LeaderApi struct {
	leaderRepotitory  p2p.LeaderRepository
	eventRepository   midgard.Repository
	messageDispatcher p2p.MessageDispatcher
	myInfo            *p2p.Node
}

type Publisher func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

func NewLeaderApi(eventRepository midgard.Repository, messageDispatcher p2p.MessageDispatcher, myInfo *p2p.Node) *LeaderApi {
	return &LeaderApi{
		eventRepository:   eventRepository,
		messageDispatcher: messageDispatcher,
		myInfo:            myInfo,
	}
}


func (leaderApi *LeaderApi) UpdateLeader(leader p2p.Leader) error {

	if leader.LeaderId.Id == "" {
		return errors.New("empty leader id purposed")
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
	}
}

func (leaderApi *LeaderApi) DeliverLeaderInfo(nodeId p2p.NodeId) error{
	if nodeId.Id == ""{
		return errors.New("empty node id purposed")
	}
	leader:= leaderApi.leaderRepotitory.GetLeader()
	leaderApi.messageDispatcher.DeliverLeaderInfo(nodeId, *leader)
}