package api

import (
	"github.com/it-chain/midgard"
	"github.com/it-chain/it-chain-Engine/blockchain"
)

type NodeApi struct {
	eventRepository midgard.Repository
	publisherId string
}

func NewNodeApi(eventRepository midgard.Repository, publisherId string) NodeApi {
	return NodeApi{
		publisherId: publisherId,
		eventRepository: eventRepository,
	}
}

// TODO
func (api *NodeApi) UpdateNode(node blockchain.Node) error {
	return nil
}
// TODO:
func (api *NodeApi) AddNode(node blockchain.Node) error {
	return nil
}
// TODO:
func (api *NodeApi) DeleteNode(node blockchain.Node) error {
	return nil
}
