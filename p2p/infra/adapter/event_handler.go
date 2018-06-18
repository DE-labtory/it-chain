package adapter

import (
	"errors"
	"log"

	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/it-chain-Engine/p2p/api"
)

var ErrEmptyAddress = errors.New("empty address proposed")

type EventHandler struct {
	nodeApi *api.NodeApi
}

func NewEventHandler(nodeApi *api.NodeApi) *EventHandler {
	return &EventHandler{
		nodeApi: nodeApi,
	}
}

func (n *EventHandler) HandleConnCreatedEvent(event p2p.ConnectionCreatedEvent) {

	if event.ID == "" || event.Address == "" {
		return
	}

	node := *p2p.NewNode(event.Address, p2p.NodeId{Id: event.ID})
	err := n.nodeApi.AddNode(node)

	if err != nil {
		log.Println(err)
	}
}

//todo conn disconnect event 구현
func (n *EventHandler) HandleConnDisconnectedEvent(event p2p.ConnectionDisconnectedEvent) {

	if event.ID == "" {
		return
	}

	err := n.nodeApi.DeleteNode(p2p.NodeId{Id: event.ID})

	if err != nil {
		log.Println(err)
	}
}

type RepositoryProjector struct {
	NodeRepository   p2p.NodeRepository
	LeaderRepository p2p.LeaderRepository
}

func NewRepositoryProjector(nodeRepository p2p.NodeRepository, leaderRepository p2p.LeaderRepository) *RepositoryProjector {
	return &RepositoryProjector{
		NodeRepository:   nodeRepository,
		LeaderRepository: leaderRepository,
	}
}

//save Leader when LeaderReceivedEvent Detected
func (projector *RepositoryProjector) HandleLeaderUpdatedEvent(event p2p.LeaderUpdatedEvent) error {

	if event.ID == "" {
		return ErrEmptyNodeId
	}

	leader := p2p.Leader{
		LeaderId: p2p.LeaderId{Id: event.ID},
	}

	projector.LeaderRepository.SetLeader(leader)
	return nil
}

func (projector *RepositoryProjector) HandlerNodeCreatedEvent(event p2p.NodeCreatedEvent) error {

	if event.ID == "" {
		return ErrEmptyNodeId
	}

	if event.IpAddress == "" {
		return ErrEmptyAddress
	}

	node := p2p.NewNode(event.IpAddress, p2p.NodeId{Id: event.ID})
	projector.NodeRepository.Save(*node)

	return nil
}
