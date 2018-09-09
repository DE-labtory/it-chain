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

package api_gateway

import (
	"sync"

	"errors"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/p2p"
)

var ErrPeerExists = errors.New("peer already exists")

type PeerQueryApi struct {
	mux            sync.Mutex
	peerRepository p2p.PeerRepository
}

func NewPeerQueryApi(repository p2p.PeerRepository) PeerQueryApi {
	return PeerQueryApi{
		mux:            sync.Mutex{},
		peerRepository: repository,
	}
}

func (pqa *PeerQueryApi) GetPLTable() (p2p.PLTable, error) {

	return pqa.peerRepository.GetPLTable()
}

func (p *PeerQueryApi) GetPeerTable() (map[string]struct {
	ID        string
	IpAddress string
}, error) {

	table, _ := p.peerRepository.GetPeerTable()

	extracted := make(map[string]struct {
		ID        string
		IpAddress string
	})

	for _, peer := range table {
		extracted[peer.PeerId.Id] = struct {
			ID        string
			IpAddress string
		}{ID: string(peer.PeerId.Id), IpAddress: peer.IpAddress}
	}

	return extracted, nil
}

func (pqa *PeerQueryApi) GetPeerList() ([]p2p.Peer, error) {
	pTable, err := pqa.peerRepository.GetPLTable()

	if err != nil {
		return nil, err
	}

	peerList := make([]p2p.Peer, 0)

	for _, peer := range pTable.PeerTable {
		peerList = append(peerList, peer)
	}

	return peerList, nil
}

func (pqa *PeerQueryApi) GetLeader() (p2p.Leader, error) {

	return pqa.peerRepository.GetLeader()
}

func (pqa *PeerQueryApi) FindPeerById(peerId p2p.PeerId) (p2p.Peer, error) {

	return pqa.peerRepository.FindPeerById(peerId)
}

func (pqa *PeerQueryApi) FindPeerByAddress(ipAddress string) (p2p.Peer, error) {

	return pqa.peerRepository.FindPeerByAddress(ipAddress)
}

type P2PEventHandler struct {
	peerRepository p2p.PeerRepository
}

func (peh *P2PEventHandler) PeerCreatedEventHandler(event event.PeerCreated) error {

	peer := p2p.Peer{
		PeerId: p2p.PeerId{
			Id: event.PeerId,
		},
		IpAddress: event.IpAddress,
	}

	peh.peerRepository.Save(peer)

	return nil
}

func (peh *P2PEventHandler) PeerDeletedEventHandler(event event.PeerCreated) error {

	peh.peerRepository.Remove(event.PeerId)

	return nil
}

func (peh *P2PEventHandler) HandleLeaderUpdatedEvent(event event.PeerCreated) error {

	leader := p2p.Leader{
		LeaderId: p2p.LeaderId{
			Id: event.PeerId,
		},
	}

	peh.peerRepository.SetLeader(leader)

	return nil

}
