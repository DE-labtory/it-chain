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
	Leader
	RepresentativeTable map[string]*Representative
	peerQueryApi        PeerQueryApi
	mux                 sync.Mutex
}

type Leader struct {
	LeaderId LeaderId
}

type LeaderId struct {
	Id string
}

func (lid LeaderId) ToString() string {
	return string(lid.Id)
}

func (l Leader) GetID() string {
	return l.LeaderId.ToString()
}

type Representative struct {
	ID RepresentativeID
}

type RepresentativeID string

func (r Representative) GetID() string {
	return string(r.ID)
}

func NewRepresentative(ID string) *Representative {
	return &Representative{ID: RepresentativeID(ID)}
}

// refresh representatives
func (p *Parliament) RefreshRepresentatives() error {
	p.mux.Lock()
	defer p.mux.Unlock()

	// get peer table from peer query api
	pt, err := p.peerQueryApi.GetPeerTable()

	// extract representatives from peer table
	for id, peer := range pt {

		p.RepresentativeTable[id] = &Representative{
			ID: RepresentativeID(peer.ID),
		}
	}

	if err != nil {
		return err
	}

	return nil
}
