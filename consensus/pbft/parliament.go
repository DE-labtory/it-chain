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
	"errors"

	"github.com/it-chain/iLogger"
)

var ErrRepresentativeDoesNotExist = errors.New("Representative does not exist")
var ErrRepresentativeAlreadyExist = errors.New("Representative already exist")

type Parliament struct {
	Leader          Leader
	Representatives map[string]Representative
}

func NewParliament() Parliament {
	return Parliament{
		Leader:          Leader{},
		Representatives: make(map[string]Representative),
	}
}

func (p Parliament) FindRepresentativeByID(representativeId string) (Representative, error) {
	r, ok := p.Representatives[representativeId]
	if !ok {
		return Representative{}, ErrRepresentativeDoesNotExist
	}

	return r, nil
}

func (p Parliament) GetRepresentatives() []Representative {
	representativeList := make([]Representative, 0)
	for _, representative := range p.Representatives {
		representativeList = append(representativeList, representative)
	}

	return representativeList
}

func (p *Parliament) AddRepresentative(representative Representative) error {
	_, ok := p.Representatives[representative.ID]
	if ok {
		return ErrRepresentativeAlreadyExist
	}

	p.Representatives[representative.ID] = representative
	return nil
}

func (p *Parliament) SetLeader(representativeId string) error {

	_, ok := p.Representatives[representativeId]
	if !ok {
		return ErrRepresentativeDoesNotExist
	}

	p.Leader = Leader{
		LeaderId: representativeId,
	}

	iLogger.Infof(nil, "[PBFT] set leader with id: %s", p.GetLeader().LeaderId)
	return nil
}

func (p *Parliament) IsNeedConsensus() bool {
	if len(p.Representatives) >= 4 {
		return true
	}
	return false
}

func (p Parliament) GetLeader() Leader {
	return p.Leader
}

func (p *Parliament) RemoveLeader() {
	p.Leader = Leader{LeaderId: ""}
}

func (p *Parliament) RemoveRepresentative(representativeId string) {
	delete(p.Representatives, representativeId)
}

type ParliamentRepository interface {
	Save(parliament Parliament)
	Load() Parliament
}

type Leader struct {
	LeaderId string
}

func (l Leader) GetID() string {
	return l.LeaderId
}

type Representative struct {
	ID string
}

func (r Representative) GetID() string {
	return string(r.ID)
}

func NewRepresentative(ID string) Representative {
	return Representative{ID: ID}
}
