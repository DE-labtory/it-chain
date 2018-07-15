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

	err := p2p.NewPeer(peer.IpAddress, peer.PeerId)

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

