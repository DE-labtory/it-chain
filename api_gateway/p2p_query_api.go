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
	peerRepository *PeerRepository
}

func NewPeerQueryApi(repository *PeerRepository) PeerQueryApi {
	return PeerQueryApi{
		mux:            sync.Mutex{},
		peerRepository: repository,
	}
}

func (pqa *PeerQueryApi) GetPLTable() (p2p.PLTable, error) {

	return pqa.peerRepository.GetPLTable()
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

type PeerRepository struct {
	mux     sync.RWMutex
	pLTable p2p.PLTable
}

func NewPeerReopository() PeerRepository {
	return PeerRepository{
		mux:     sync.RWMutex{},
		pLTable: *p2p.NewPLTable(p2p.Leader{p2p.LeaderId{""}}, make(map[string]p2p.Peer)),
	}
}

func (pltrepo *PeerRepository) GetPLTable() (p2p.PLTable, error) {

	pltrepo.mux.Lock()
	defer pltrepo.mux.Unlock()

	return pltrepo.pLTable, nil
}

func (pltrepo *PeerRepository) GetLeader() (p2p.Leader, error) {

	return pltrepo.pLTable.Leader, nil
}

func (pltrepo *PeerRepository) FindPeerById(peerId p2p.PeerId) (p2p.Peer, error) {

	pltrepo.mux.Lock()
	defer pltrepo.mux.Unlock()
	v, exist := pltrepo.pLTable.PeerTable[peerId.Id]

	if peerId.Id == "" {
		return v, p2p.ErrEmptyPeerId
	}
	//no matching id
	if !exist {
		return v, p2p.ErrNoMatchingPeerId
	}

	return v, nil
}

func (pltrepo *PeerRepository) FindPeerByAddress(ipAddress string) (p2p.Peer, error) {

	pltrepo.mux.Lock()
	defer pltrepo.mux.Unlock()

	for _, peer := range pltrepo.pLTable.PeerTable {

		if peer.IpAddress == ipAddress {
			return peer, nil
		}
	}

	return p2p.Peer{}, nil
}

func (pltrepo *PeerRepository) Save(peer p2p.Peer) error {

	pltrepo.mux.Lock()
	defer pltrepo.mux.Unlock()

	_, exist := pltrepo.pLTable.PeerTable[peer.PeerId.Id]

	if exist {
		return ErrPeerExists
	}

	pltrepo.pLTable.PeerTable[peer.PeerId.Id] = peer

	return nil
}

func (pltrepo *PeerRepository) SetLeader(peer p2p.Peer) error {

	pltrepo.mux.Lock()
	defer pltrepo.mux.Unlock()

	leader := p2p.Leader{
		LeaderId: p2p.LeaderId{
			Id: peer.PeerId.Id,
		},
	}

	pltrepo.pLTable.Leader = leader

	return nil
}

func (pltrepo *PeerRepository) Remove(id string) error {

	pltrepo.mux.Lock()
	defer pltrepo.mux.Unlock()

	delete(pltrepo.pLTable.PeerTable, id)

	return nil
}

type P2PEventHandler struct {
	peerRepository PeerRepository
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

	peer := p2p.Peer{
		PeerId: p2p.PeerId{
			Id: event.PeerId,
		},
	}

	peh.peerRepository.SetLeader(peer)

	return nil

}
