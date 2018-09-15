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

package pbft

import (
	"sync"
)

type Parliament struct {
	Leader              *Leader
	RepresentativeTable map[string]*Representative
	peerQueryApi        PeerQueryApi
	mux                 sync.Mutex
}

type Leader struct {
	LeaderId string
}

func (l Leader) GetID() string {
	return l.LeaderId
}

type Representative struct {
	ID        string
	IpAddress string
}

func (r Representative) GetID() string {
	return string(r.ID)
}

func NewRepresentative(ID string) *Representative {
	return &Representative{ID: ID}
}

func (p *Parliament) SetLeader(representative *Representative) error {
	p.mux.Lock()
	defer p.mux.Unlock()

	p.Leader = &Leader{
		LeaderId: representative.ID,
	}

	return nil

}

// build parliament
func (p *Parliament) Build() error {
	p.mux.Lock()
	defer p.mux.Unlock()

	// get peer table from peer query api
	pt, err := p.peerQueryApi.GetPeerTable()

	// extract representatives from peer table
	for id, peer := range pt {

		p.RepresentativeTable[id] = &Representative{
			ID: peer.ID,
		}
	}

	if err != nil {
		return err
	}

	return nil
}

func (p *Parliament) FindRepresentativeByIpAddress(ipAddress string) *Representative {
	for _, rep := range p.RepresentativeTable {

		if rep.IpAddress == ipAddress {

			return rep
		}

	}

	return nil
}
