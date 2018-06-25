package api

import "github.com/it-chain/midgard"

type NodeApi struct {
	eventRepository *midgard.Repository
	publisherId string
}

func NewNodeApi(eventRepository *midgard.Repository, publisherId string) NodeApi {
	return NodeApi{
		publisherId: publisherId,
		eventRepository: eventRepository,
	}
}

func (n NodeApi) AddNode() {

}
