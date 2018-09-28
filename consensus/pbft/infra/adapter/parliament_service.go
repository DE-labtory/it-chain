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
	"sync"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/consensus/pbft"
)

type ParliamentService struct {
	Parliament   *pbft.Parliament
	peerQueryApi *api_gateway.PeerQueryApi
	mux          sync.Mutex
}

func NewParliamentService(parliament *pbft.Parliament, peerQueryApi *api_gateway.PeerQueryApi) *ParliamentService {
	return &ParliamentService{
		Parliament:   parliament,
		peerQueryApi: peerQueryApi,
		mux:          sync.Mutex{},
	}
}

func (ps *ParliamentService) RequestLeader() (pbft.MemberID, error) {
	l, err := ps.peerQueryApi.GetLeader()

	if err != nil {
		return "", err
	}

	return pbft.MemberID(l.ID), nil
}

func (ps *ParliamentService) RequestPeerList() ([]pbft.MemberID, error) {

	pl := ps.peerQueryApi.GetAllPeerList()
	peerList := make([]pbft.MemberID, 0)

	for _, p := range pl {
		peerList = append(peerList, pbft.MemberID(p.ID))
	}

	return peerList, nil
}

func (p *ParliamentService) IsNeedConsensus() bool {
	peerList, err := p.RequestPeerList()

	if err != nil {
		return false
	}

	numOfMember := 0
	numOfMember = numOfMember + len(peerList)

	if numOfMember >= 4 {
		return true
	}

	return false
}

// build parliament
func (p *ParliamentService) Build() error {
	p.mux.Lock()
	defer p.mux.Unlock()

	// get peer table from peer query api
	peerList := p.peerQueryApi.GetAllPeerList()

	// extract representatives from peer table
	for _, peer := range peerList {
		p.Parliament.RepresentativeTable[peer.ID] = &pbft.Representative{
			ID:        peer.ID,
			IpAddress: peer.GrpcGatewayAddress,
		}
	}

	return nil
}

func (p *ParliamentService) SetLeader(representative *pbft.Representative) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.Parliament.Leader = &pbft.Leader{
		LeaderId: representative.ID,
	}

	return nil
}

func (p *ParliamentService) GetRepresentativeById(id string) *pbft.Representative {
	p.mux.Lock()
	defer p.mux.Unlock()

	return p.Parliament.RepresentativeTable[id]
}

func (p *ParliamentService) GetRepresentativeTable() map[string]*pbft.Representative {
	p.mux.Lock()
	defer p.mux.Unlock()

	return p.Parliament.RepresentativeTable
}

func (p *ParliamentService) GetParliament() *pbft.Parliament {
	p.mux.Lock()
	defer p.mux.Unlock()

	return p.Parliament
}

func (p *ParliamentService) GetLeader() *pbft.Leader {
	p.mux.Lock()
	defer p.mux.Unlock()

	return p.Parliament.Leader
}

func (p *ParliamentService) FindRepresentativeByIpAddress(ipAddress string) *pbft.Representative {
	for _, rep := range p.Parliament.RepresentativeTable {

		if rep.IpAddress == ipAddress {

			return rep
		}

	}
	return nil
}
