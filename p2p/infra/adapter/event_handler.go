package adapter

import (
	"errors"
	"log"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrEmptyAddress = errors.New("empty address proposed")
var ErrPeerApi = errors.New("problem in node api")
type EventHandlerPeerApi interface{
	AddPeer(node p2p.Peer) error
	DeletePeer(id p2p.PeerId) error
}
type EventHandler struct {
	nodeApi EventHandlerPeerApi
}

func NewEventHandler(nodeApi EventHandlerPeerApi) *EventHandler {
	return &EventHandler{
		nodeApi: nodeApi,
	}
}


func (n *EventHandler) HandleConnCreatedEvent(event p2p.ConnectionCreatedEvent) error {
	if event.ID == "" {
		return ErrEmptyPeerId
	}

	if event.Address == "" {
		return ErrEmptyAddress
	}

	node := *p2p.NewPeer(event.Address, p2p.PeerId{Id: event.ID})
	err := n.nodeApi.AddPeer(node)

	if err != nil {
		return ErrPeerApi
	}

	return nil
}

//todo conn disconnect event 구현
func (n *EventHandler) HandleConnDisconnectedEvent(event p2p.ConnectionDisconnectedEvent) error {

	if event.ID == "" {
		return ErrEmptyPeerId
	}

	err := n.nodeApi.DeletePeer(p2p.PeerId{Id: event.ID})

	if err != nil {
		log.Println(err)
	}

	return nil
}

type WriteOnlyPeerRepository interface{
	Save(data p2p.Peer) error
}
type WriteOnlyLeaderRepository interface {
	SetLeader(leader p2p.Leader)
}
type RepositoryProjector struct {
	nodeRepository   WriteOnlyPeerRepository
	leaderRepository WriteOnlyLeaderRepository
}

func NewRepositoryProjector(nodeRepository WriteOnlyPeerRepository, leaderRepository WriteOnlyLeaderRepository) *RepositoryProjector {
	return &RepositoryProjector{
		nodeRepository:   nodeRepository,
		leaderRepository: leaderRepository,
	}
}

//save Leader when LeaderReceivedEvent Detected
func (projector *RepositoryProjector) HandleLeaderUpdatedEvent(event p2p.LeaderUpdatedEvent) error {

	if event.ID == "" {
		return ErrEmptyPeerId
	}

	leader := p2p.Leader{
		LeaderId: p2p.LeaderId{Id: event.ID},
	}

	projector.leaderRepository.SetLeader(leader)
	return nil
}

func (projector *RepositoryProjector) HandlerPeerCreatedEvent(event p2p.PeerCreatedEvent) error {

	if event.ID == "" {
		return ErrEmptyPeerId
	}

	if event.IpAddress == "" {
		return ErrEmptyAddress
	}

	node := p2p.NewPeer(event.IpAddress, p2p.PeerId{Id: event.ID})
	projector.nodeRepository.Save(*node)

	return nil
}
