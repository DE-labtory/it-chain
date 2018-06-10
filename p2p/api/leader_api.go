package api

import (
	"log"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
)

type LeaderApi struct {
	leaderRepotitory  p2p.LeaderRepository
	eventRepository   midgard.Repository
	messageDispatcher p2p.MessageDispatcher
	myInfo            *p2p.Node
}


func NewLeaderApi(eventRepository midgard.Repository, messageDispatcher p2p.MessageDispatcher, myInfo *p2p.Node) *LeaderApi {
	return &LeaderApi{
		eventRepository:   eventRepository,
		messageDispatcher: messageDispatcher,
		myInfo:            myInfo,
	}
}


func (leaderApi *LeaderApi) UpdateLeader(leader p2p.Node) {

	if leader.NodeId == "" {
		return
	}

	events := make([]midgard.Event, 0)
	leaderUpdatedEvent := p2p.LeaderUpdatedEvent{
		EventModel: midgard.EventModel{
			ID:   leader.NodeId.ToString(),
			Type: "leader.update",
		},
	}

	events = append(events, leaderUpdatedEvent)
	err := leaderApi.eventRepository.Save(leaderUpdatedEvent.GetID(), events...)

	if err != nil {
		log.Println(err.Error())
	}
}

func (leaderApi *LeaderApi) DeliverLeaderInfo(nodeId p2p.NodeId){
	leader:= leaderApi.leaderRepotitory.GetLeader()
	leaderApi.messageDispatcher.DeliverLeaderInfo(nodeId, *leader)
}