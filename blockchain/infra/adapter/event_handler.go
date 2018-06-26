package adapter

import "github.com/it-chain/it-chain-Engine/blockchain"

type EventHandlerNodeApi interface {
	AddNode(node blockchain.Node) error
	DeleteNode(node blockchain.Node) error
}

type EventHandler struct {
	nodeApi EventHandlerNodeApi
}

func NewEventHandler(nodeApi EventHandlerNodeApi) *EventHandler{
	return &EventHandler{
		nodeApi: nodeApi,
	}
}

/// Check 단계에서 임의의 노드를 선정하기 위해 노드를 저장한다.
func (eh *EventHandler) HandleNodeCreatedEvent(event blockchain.NodeCreatedEvent) error {
	return nil
}

func (eh *EventHandler) HandleNodeDeletedEvent(event blockchain.NodeDeletedEvent) error {
	return nil
}

