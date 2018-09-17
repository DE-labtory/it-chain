/*
 * Copyright 2018 It-chain
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package adapter

import (
	"errors"
	"log"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/api"
)

var ErrPeerApi = errors.New("problem in peer api")

type CommunicationApi interface {
	DialToUnConnectedNode(peerTable map[string]p2p.Peer) error
	DeliverPLTable(connectionId string) error
}

type EventHandler struct {
	communicationApi CommunicationApi // api.CommunicationApi
	peerApi          api.PeerApi
}

func NewEventHandler(communicationApi CommunicationApi, peerApi api.PeerApi) EventHandler {

	return EventHandler{
		communicationApi: communicationApi,
		peerApi:          peerApi,
	}
}

//handler connection created event
func (eh *EventHandler) HandleConnCreatedEvent(event event.ConnectionCreated) error {

	//1. addPeer
	peer := p2p.Peer{
		PeerId: p2p.PeerId{
			Id: event.ConnectionID,
		},
		IpAddress: event.Address,
	}

	err := eh.peerApi.Save(peer)

	if err != nil {
		return err
	}

	//2. send peer table
	eh.communicationApi.DeliverPLTable(event.ConnectionID)

	return nil
}

//todo deleted peer if disconnected peer is leader
func (eh *EventHandler) HandleConnDisconnectedEvent(event event.ConnectionClosed) error {

	if event.ConnectionID == "" {
		return ErrEmptyPeerId
	}

	err := eh.peerApi.Remove(p2p.PeerId{Id: event.ConnectionID})

	if err != nil {
		log.Println(err)
	}

	return nil
}
