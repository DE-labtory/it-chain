package adapter

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"errors"
	"log"
	"golang.org/x/tools/go/gcimporter15/testdata"
	"github.com/it-chain/midgard"
)


var ErrEmptyEventId = errors.New("empty event id proposed.")
var ErrEmptyIpAddress = errors.New("empty ip address proposed.")
var ErrEmptyNodeId = errors.New("empty node id proposed.")
var ErrNodeApi = errors.New("problem in node api")

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
	eventID := event.GetID()

	if eventID == "" {
		return ErrEmptyEventId
	}

	node := event.Node

	if err := isValidNode(node); err != nil {
		return err
	}

	err := eh.nodeApi.AddNode(node)

	if err != nil {
		return ErrNodeApi
	}

	return nil
}

func (eh *EventHandler) HandleNodeDeletedEvent(event blockchain.NodeDeletedEvent) error {
	return nil
}

func isValidNode(node blockchain.Node) error {
	if node.IpAddress == "" {
		return ErrEmptyIpAddress
	}

	if node.NodeId.Id == "" {
		return ErrEmptyNodeId
	}
}

