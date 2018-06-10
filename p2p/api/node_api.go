package api

import (
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/midgard"
	"github.com/it-chain/it-chain-Engine/gateway"
)

type NodeApi struct {
	nodeRepository *p2p.NodeRepository
	leaderRepository *p2p.LeaderRepository
	eventRepository midgard.Repository
	messageDispatcher *p2p.MessageDispatcher
}

func NewNodeApi(nodeRepository *p2p.NodeRepository, leaderRepository *p2p.LeaderRepository, eventRepository *midgard.Repository, messageDispatcher *p2p.MessageDispatcher) *NodeApi{
	return &NodeApi{
		nodeRepository : nodeRepository,
		leaderRepository: leaderRepository,
		eventRepository: eventRepository,
		messageDispatcher: messageDispatcher,
	}
}

func (nodeApi *NodeApi) UpdateNodeList(command gateway.MessageReceiveCommand) {
	if command.GetID() ==""{
		return
	}

	id := command.GetID()

	nodeList := make([]p2p.Node,0)
	err := common.Deserialize(command.Data, nodeList)

	if err != nil{
		err.Error()
	}

	event := p2p.NodeListUpdatedEvent{
		EventModel: midgard.EventModel{
			ID:id,
			Type:"Node",
		},
		NodeList:nodeList,
	}

	nodeApi.messageDispatcher.publisher.Publish("event", "node.update", event)
}

func (nodeApi *NodeApi) DeliverNodeList(command gateway.MessageDeliverCommand){
	deliverCommand := Me
}
