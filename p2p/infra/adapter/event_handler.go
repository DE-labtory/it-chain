package adapter

import (
	"errors"
	"log"

	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/api"
)

var ErrPeerApi = errors.New("problem in peer api")

type EventHandler struct {
	communicationApi api.ICommunicationApi
}

func NewEventHandler(communicationApi api.ICommunicationApi) *EventHandler {

	return &EventHandler{
		communicationApi: communicationApi,
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
	eh.communicationApi.DeliverPLTable(event.ID)

	return nil
}

//todo deleted peer if disconnected peer is leader
func (eh *EventHandler) HandleConnDisconnectedEvent(event p2p.ConnectionDisconnectedEvent) error {

	if event.ID == "" {
		return ErrEmptyPeerId
	}

	err := p2p.DeletePeer(p2p.PeerId{Id: event.ID})

	if err != nil {
		log.Println(err)
	}

	return nil
}

