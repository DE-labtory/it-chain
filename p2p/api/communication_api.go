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

package api

import "github.com/it-chain/engine/p2p"

type ICommunicationApi interface {
	DialToUnConnectedNode(peerTable map[string]p2p.Peer) error
	DeliverPLTable(connectionId string) error
}

type CommunicationApi struct {
	peerQueryService     p2p.PeerQueryService
	communicationService p2p.ICommunicationService
}

func NewCommunicationApi(peerQueryService p2p.PeerQueryService, communicationService p2p.ICommunicationService) *CommunicationApi {
	return &CommunicationApi{
		peerQueryService:     peerQueryService,
		communicationService: communicationService,
	}
}

func (ca *CommunicationApi) DialToUnConnectedNode(peerTable map[string]p2p.Peer) error {

	//1. find unconnected peer
	//2. dial to unconnected peer
	for _, peer := range peerTable {

		//err is nil if there is matching peer
		peer, err := ca.peerQueryService.FindPeerById(peer.PeerId)

		//dial if no peer matching peer id
		if err != nil {
			ca.communicationService.Dial(peer.IpAddress)
		}
	}

	return nil
}

//Deliver Peer leader table that consists of peerList and leader
func (ca *CommunicationApi) DeliverPLTable(connectionId string) error {

	//1. get peer table
	peerTable, _ := ca.peerQueryService.GetPLTable()

	//2. deliver peer table
	ca.communicationService.DeliverPLTable(connectionId, peerTable)

	return nil
}
