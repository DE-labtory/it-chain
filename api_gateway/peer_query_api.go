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
	"errors"
	"sync"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/iLogger"
)

var ErrPeerExists = errors.New("peer already exists")
var ErrPeerDoesNotExists = errors.New("peer does not exist")
var ErrLeaderDoesNotExists = errors.New("leader does not exist")

type Role string

const (
	Leader Role = "Leader"
	Member Role = "Member"
)

type Peer struct {
	ID                 string
	GrpcGatewayAddress string
	ApiGatewayAddress  string
	Role               Role
}

type PeerQueryApi struct {
	peerRepository *PeerRepository
}

func NewPeerQueryApi(peerRepository *PeerRepository) *PeerQueryApi {
	return &PeerQueryApi{
		peerRepository: peerRepository,
	}
}

func (p PeerQueryApi) GetAllPeerList() []Peer {
	return p.peerRepository.FindAll()
}

func (p PeerQueryApi) GetPeerByID(peerID string) (Peer, error) {
	return p.peerRepository.FindById(peerID)
}

func (p PeerQueryApi) GetLeader() (Peer, error) {
	return p.peerRepository.GetLeader()
}

func (p PeerQueryApi) GetNetwork() (Network, error) {
	leader, err := p.peerRepository.GetLeader()
	if err != nil {
		return Network{}, err
	}

	return Network{
		Leader:  leader,
		Members: p.peerRepository.FindAll(),
	}, nil
}

type Network struct {
	Leader  Peer
	Members []Peer
}

type PeerRepository struct {
	sync.RWMutex
	peers map[string]Peer
}

func NewPeerRepository() *PeerRepository {
	return &PeerRepository{
		peers:   make(map[string]Peer),
		RWMutex: sync.RWMutex{},
	}
}

func (p *PeerRepository) Save(peer Peer) error {

	p.Lock()
	defer p.Unlock()

	_, ok := p.peers[peer.ID]
	if ok {
		return ErrPeerExists
	}

	p.peers[peer.ID] = peer
	return nil
}

func (p *PeerRepository) Remove(ID string) {

	p.Lock()
	defer p.Unlock()

	delete(p.peers, ID)
}

func (p *PeerRepository) FindById(ID string) (Peer, error) {
	peer, ok := p.peers[ID]
	if !ok {
		return Peer{}, ErrPeerDoesNotExists
	}

	return peer, nil
}

func (p *PeerRepository) FindAll() []Peer {
	peerList := make([]Peer, 0)

	for _, peer := range p.peers {
		peerList = append(peerList, peer)
	}

	return peerList
}

func (p *PeerRepository) SetLeader(ID string) error {
	p.Lock()
	defer p.Unlock()

	newLeader, ok := p.peers[ID]
	if !ok {
		return ErrPeerDoesNotExists
	}

	for _, peer := range p.peers {
		peer.Role = Member
	}

	newLeader.Role = Leader
	p.peers[ID] = newLeader

	return nil
}

func (p *PeerRepository) GetLeader() (Peer, error) {
	p.Lock()
	defer p.Unlock()

	for _, peer := range p.peers {
		if peer.Role == Leader {
			return peer, nil
		}
	}

	return Peer{}, ErrLeaderDoesNotExists
}

type ConnectionEventHandler struct {
	peerRepository *PeerRepository
}

func NewConnectionEventListener(peerRepository *PeerRepository) *ConnectionEventHandler {
	return &ConnectionEventHandler{
		peerRepository: peerRepository,
	}
}

func (c *ConnectionEventHandler) HandleConnectionCreatedEvent(event event.ConnectionCreated) {

	peer := Peer{
		ID:                 event.ConnectionID,
		GrpcGatewayAddress: event.GrpcGatewayAddress,
		ApiGatewayAddress:  event.ApiGatewayAddress,
		Role:               Member,
	}

	if err := c.peerRepository.Save(peer); err != nil {
		iLogger.Errorf(nil, "[Api-gateway] Fail to save peer - Err:[%s]", err.Error())
	}
}

func (c *ConnectionEventHandler) HandleConnectionClosedEvent(event event.ConnectionClosed) {
	c.peerRepository.Remove(event.ConnectionID)
}

type LeaderUpdateEventListener struct {
	peerRepository *PeerRepository
}

func NewLeaderUpdateEventListener(peerRepository *PeerRepository) *LeaderUpdateEventListener {
	return &LeaderUpdateEventListener{
		peerRepository: peerRepository,
	}
}

func (l *LeaderUpdateEventListener) HandleLeaderUpdatedEvent(event event.LeaderUpdated) {
	if err := l.peerRepository.SetLeader(event.LeaderId); err != nil {
		iLogger.Errorf(nil, "[Api-gateway] Fail to set leader - Err:[%s]", err.Error())
	}
}
