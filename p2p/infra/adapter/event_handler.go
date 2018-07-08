package adapter

import (
	"errors"
	"log"

	"github.com/it-chain/it-chain-Engine/p2p"
)

var ErrEmptyAddress = errors.New("empty address proposed")
var ErrPeerApi = errors.New("problem in peer api")

type EventHandlerPeerApi interface {
	AddPeer(peer p2p.Peer) error
	DeletePeer(id p2p.PeerId) error
	DeliverPLTable(connectionId string) error
}
type EventHandler struct {
	peerApi EventHandlerPeerApi
}

func NewEventHandler(peerApi EventHandlerPeerApi) *EventHandler {
	return &EventHandler{
		peerApi: peerApi,
	}
}

//handler connection created event
func (eh *EventHandler) HandleConnCreatedEvent(event p2p.ConnectionCreatedEvent) error {

	//1. addPeer
	peer := p2p.Peer{
		PeerId:p2p.PeerId{
			Id:event.ID,
		},
		IpAddress:event.Address,

	}

	err := eh.peerApi.AddPeer(peer)

	if err != nil {
		return err
	}

	//2. send peer table
	eh.peerApi.DeliverPLTable(event.ID)

	return nil
}

//todo deleted peer if disconnected peer is leader
func (eh *EventHandler) HandleConnDisconnectedEvent(event p2p.ConnectionDisconnectedEvent) error {

	if event.ID == "" {
		return ErrEmptyPeerId
	}

	err := eh.peerApi.DeletePeer(p2p.PeerId{Id: event.ID})

	if err != nil {
		log.Println(err)
	}

	return nil
}

type WriteOnlyPeerRepository interface {
	Save(data p2p.Peer) error
}

type WriteOnlyLeaderRepository interface {
	SetLeader(leader p2p.Leader)
}

type RepositoryProjector struct {
	peerRepository   WriteOnlyPeerRepository
	leaderRepository WriteOnlyLeaderRepository
}

func NewRepositoryProjector(peerRepository WriteOnlyPeerRepository, leaderRepository WriteOnlyLeaderRepository) *RepositoryProjector {

	return &RepositoryProjector{
		peerRepository:   peerRepository,
		leaderRepository: leaderRepository,
	}
}

//save Leader when LeaderReceivedEvent Detected, publish updated info to network
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

	peer := p2p.NewPeer(event.IpAddress, p2p.PeerId{Id: event.ID})

	projector.peerRepository.Save(peer)

	return nil
}
