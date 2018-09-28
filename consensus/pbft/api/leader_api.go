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
)

var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyConnectionId = errors.New("empty connection id proposed")
var ErrNoMatchingPeerWithIpAddress = errors.New("no matching peer with ip address")

type LeaderApi struct {
	parliamentService pbft.ParliamentService
	eventService      common.EventService
}

func NewLeaderApi(ps pbft.ParliamentService, eventService common.EventService) *LeaderApi {

	return &LeaderApi{
		parliamentService: ps,
		eventService:      eventService,
	}
}

func (la *LeaderApi) UpdateLeader(nodeId string) error {
	//1. loop peer list and find specific address
	//2. update specific peer as leader
	rep := la.parliamentService.GetRepresentativeById(nodeId)
	if rep == nil {
		return ErrNoMatchingPeerWithIpAddress
	}

	la.parliamentService.SetLeader(&pbft.Representative{
		ID: rep.ID,
	})

	return la.eventService.Publish("leader.updated", event.LeaderUpdated{
		LeaderId: rep.ID,
	})

}

func (la *LeaderApi) GetParliament() *pbft.Parliament {
	return la.parliamentService.GetParliament()
}

func (la *LeaderApi) GetLeader() *pbft.Leader {
	return la.parliamentService.GetLeader()
}
