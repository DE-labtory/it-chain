package api

import (
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/common"
	"log"
	"github.com/it-chain/midgard"
	"github.com/it-chain/it-chain-Engine/gateway"
)

type LeaderApi struct {
	nodeRepository p2p.NodeRepository
	leaderRepository p2p.LeaderRepository
	eventRepository midgard.Repository
	messageDispatcher p2p.MessageDispatcher
}

func NewLeaderApi(nodeRepository *p2p.NodeRepository, leaderRepository *p2p.LeaderRepository, eventRepository *midgard.Repository, messageDispatcher *p2p.MessageDispatcher) *LeaderApi{
	return &LeaderApi{
		nodeRepository : nodeRepository,
		leaderRepository: leaderRepository,
		eventRepository: eventRepository,
		messageDispatcher: messageDispatcher,
	}
}

func (leaderApi *LeaderApi) UpdateLeader(command gateway.MessageReceiveCommand) {
	if command.GetID() == ""{
		return
	}
	id := command.GetID()
	leader := &p2p.Leader{}
	err := common.Deserialize(command.Data, leader)
	if err != nil{
		panic(err)
	}

	events := make([]midgard.Event, 0)
	leaderUpdatedEvent := p2p.LeaderUpdatedEvent{
		EventModel: midgard.EventModel{
			ID: 	id,
			Type:	"Leader",
		},
		Leader: *leader,
	}

	events = append(events, leaderUpdatedEvent)
	err2 := leaderApi.eventRepository.Save(command.GetID(), events...)

	if err2 != nil {
		log.Println(err2.Error())
	}

	leaderApi.messageDispatcher.publisher.Publish("event", "leader.update", leaderUpdatedEvent)
}


