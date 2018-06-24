package adapter

import (
	"errors"
	"log"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrEmptyAddress = errors.New("empty address proposed")
var ErrNodeApi = errors.New("problem in node api")
type EventHandlerNodeApi interface{
	AddNode(node p2p.Node) error
	DeleteNode(id p2p.NodeId) error
}
type EventHandler struct {
	nodeApi EventHandlerNodeApi
}

func NewEventHandler(nodeApi EventHandlerNodeApi) *EventHandler {
	return &EventHandler{
		nodeApi: nodeApi,
	}
}


func (n *EventHandler) HandleConnCreatedEvent(event p2p.ConnectionCreatedEvent) error {
	if event.ID == "" {
		return ErrEmptyNodeId
	}

	if event.Address == "" {
		return ErrEmptyAddress
	}

	node := *p2p.NewNode(event.Address, p2p.NodeId{Id: event.ID})
	err := n.nodeApi.AddNode(node)

	if err != nil {
		return ErrNodeApi
	}

	return nil
}

//todo conn disconnect event 구현
func (n *EventHandler) HandleConnDisconnectedEvent(event p2p.ConnectionDisconnectedEvent) error {

	if event.ID == "" {
		return ErrEmptyNodeId
	}

	err := n.nodeApi.DeleteNode(p2p.NodeId{Id: event.ID})

	if err != nil {
		log.Println(err)
	}

	return nil
}

type WriteOnlyNodeRepository interface{
	Save(data p2p.Node) error
}
type WriteOnlyLeaderRepository interface {
	SetLeader(leader p2p.Leader)
}
type RepositoryProjector struct {
	nodeRepository   WriteOnlyNodeRepository
	leaderRepository WriteOnlyLeaderRepository
}

func NewRepositoryProjector(nodeRepository WriteOnlyNodeRepository, leaderRepository WriteOnlyLeaderRepository) *RepositoryProjector {
	return &RepositoryProjector{
		nodeRepository:   nodeRepository,
		leaderRepository: leaderRepository,
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

	projector.leaderRepository.SetLeader(leader)
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
	projector.nodeRepository.Save(*node)

	return nil
}
