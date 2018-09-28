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

import (
	"errors"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/iLogger"
)

var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyConnectionId = errors.New("empty connection id proposed")
var ErrNoMatchingPeerWithIpAddress = errors.New("no matching peer with ip address")

type ParliamentApi struct {
	parliamentRepository pbft.ParliamentRepository
	eventService         common.EventService
}

func NewParliamentApi(parliamentRepository pbft.ParliamentRepository, eventService common.EventService) *ParliamentApi {

	return &ParliamentApi{
		parliamentRepository: parliamentRepository,
		eventService:         eventService,
	}
}

func (p *ParliamentApi) AddRepresentative(representativeId string) {
	parliament := p.parliamentRepository.Load()
	parliament.AddRepresentative(pbft.Representative{
		ID: representativeId,
	})

	p.parliamentRepository.Save(parliament)
}

func (p *ParliamentApi) RemoveRepresentative(representativeId string) {
	parliament := p.parliamentRepository.Load()
	parliament.RemoveRepresentative(representativeId)

	p.parliamentRepository.Save(parliament)
}

func (p *ParliamentApi) UpdateLeader(nodeId string) error {
	//1. loop peer list and find specific address
	//2. update specific peer as leader

	parliament := p.parliamentRepository.Load()
	representative, err := parliament.FindRepresentativeByID(nodeId)
	iLogger.Infof(nil, "[PBFT] found representative to be leader id: %s", representative.ID)
	if err != nil {
		return ErrNoMatchingPeerWithIpAddress
	}

	if err := parliament.SetLeader(representative.ID); err != nil {
		return err
	}

	p.parliamentRepository.Save(parliament)

	return p.eventService.Publish("leader.updated", event.LeaderUpdated{LeaderId: representative.ID})
}

func (p *ParliamentApi) GetLeader() pbft.Leader {
	parliament := p.parliamentRepository.Load()
	return parliament.GetLeader()
}
