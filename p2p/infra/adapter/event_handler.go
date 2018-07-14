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
		PeerId: p2p.PeerId{
			Id: event.ID,
		},
		IpAddress: event.Address,
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


type RepositoryProjector struct {}

//save Leader when LeaderReceivedEvent Detected, publish updated info to network
func (projector *RepositoryProjector) HandleLeaderUpdatedEvent(event p2p.LeaderUpdatedEvent) error {

	if event.ID == "" {
		return ErrEmptyPeerId
	}

	peer := p2p.Peer{
		PeerId: p2p.PeerId{Id: event.ID},
	}

	p2p.UpdateLeader(peer)

	return nil

}

func (projector *RepositoryProjector) HandlerPeerCreatedEvent(event p2p.PeerCreatedEvent) error {

	if event.ID == "" {
		return ErrEmptyPeerId
	}

	if event.IpAddress == "" {
		return ErrEmptyAddress
	}

	return p2p.NewPeer(event.IpAddress, p2p.PeerId{Id: event.ID})

}
