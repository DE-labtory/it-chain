package api

import (
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
	"errors"
)

var ErrEmptyNodeList = errors.New("empty node list proposed")

type ReadOnlyNodeRepository interface {
	FindById(id p2p.NodeId) (*p2p.Node, error)
	FindAll() ([]p2p.Node, error)
}

type NodeApi struct {

	nodeRepository    ReadOnlyNodeRepository
	eventRepository   midgard.Repository
	messageDispatcher p2p.MessageDispatcher
}

func NewNodeApi(nodeRepository ReadOnlyNodeRepository, eventRepository midgard.Repository, messageDispatcher p2p.MessageDispatcher) *NodeApi {
	return &NodeApi{
		nodeRepository:    nodeRepository,
		eventRepository:   eventRepository,
		messageDispatcher: messageDispatcher,
	}
}


func (nodeApi *NodeApi) UpdateNodeList(nodeList []p2p.Node) error {

	//둘다 존재할경우 무시, existNodeList에만 존재할경우 NodeDeletedEvent, nodeList에 존재할경우 NodeCreatedEvent
	var event midgard.Event

	existNodeList, err := nodeApi.nodeRepository.FindAll()

	if err != nil {
		return err
	}

	newNodes, disconnectedNodes := p2p.GetMutuallyExclusiveNodes(nodeList, existNodeList)

	for _, node := range newNodes {

		event = p2p.NodeCreatedEvent{
			EventModel: midgard.EventModel{
				ID:   node.GetID(),
				Type: "node.created",
			},
			IpAddress: node.IpAddress,
		}

		nodeApi.eventRepository.Save(event.GetID(), event)
	}

	for _, node := range disconnectedNodes {
		event = p2p.NodeDeletedEvent{
			EventModel: midgard.EventModel{
				ID:   node.GetID(),
				Type: "node.deleted",
			},
		}

		nodeApi.eventRepository.Save(event.GetID(), event)
	}
}

func (nodeApi *NodeApi) DeliverNodeList(nodeId p2p.NodeId){
	nodeList, _ := nodeApi.nodeRepository.FindAll()
	nodeApi.messageDispatcher.DeliverNodeList(nodeId, nodeList)
}