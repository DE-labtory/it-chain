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

package mem

import (
	"errors"
	"sync"

	"github.com/it-chain/engine/p2p"
)

var ErrPeerExists = errors.New("peer already exists")

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

func (pr *PeerRepository) GetPLTable() (p2p.PLTable, error) {

	pr.mux.Lock()
	defer pr.mux.Unlock()

	return pr.pLTable, nil
}

func (pr *PeerRepository) GetLeader() (p2p.Leader, error) {

	return pr.pLTable.Leader, nil
}

func (pr *PeerRepository) FindPeerById(peerId p2p.PeerId) (p2p.Peer, error) {

	pr.mux.Lock()
	defer pr.mux.Unlock()
	v, exist := pr.pLTable.PeerTable[peerId.Id]

	if peerId.Id == "" {
		return v, p2p.ErrEmptyPeerId
	}
	//no matching id
	if !exist {
		return v, p2p.ErrNoMatchingPeerId
	}

	return v, nil
}

func (pr *PeerRepository) FindPeerByAddress(ipAddress string) (p2p.Peer, error) {

	pr.mux.Lock()
	defer pr.mux.Unlock()

	for _, peer := range pr.pLTable.PeerTable {

		if peer.IpAddress == ipAddress {
			return peer, nil
		}
	}

	return p2p.Peer{}, nil
}

func (pr *PeerRepository) Save(peer p2p.Peer) error {

	pr.mux.Lock()
	defer pr.mux.Unlock()

	_, exist := pr.pLTable.PeerTable[peer.PeerId.Id]

	if exist {
		return ErrPeerExists
	}

	pr.pLTable.PeerTable[peer.PeerId.Id] = peer

	return nil
}

func (pr *PeerRepository) SetLeader(leader p2p.Leader) error {

	pr.mux.Lock()
	defer pr.mux.Unlock()

	pr.pLTable.Leader = leader

	return nil
}

func (pr *PeerRepository) Remove(id string) error {

	pr.mux.Lock()
	defer pr.mux.Unlock()

	delete(pr.pLTable.PeerTable, id)

	return nil
}
