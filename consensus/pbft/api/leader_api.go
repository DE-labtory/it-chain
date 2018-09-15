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
	parliament   *pbft.Parliament
	eventService common.EventService
}

func NewLeaderApi(parliament *pbft.Parliament, eventService common.EventService) *LeaderApi {

	return &LeaderApi{
		parliament:   parliament,
		eventService: eventService,
	}
}

func (la *LeaderApi) UpdateLeaderWithAddress(ipAddress string) error {
	//1. loop peer list and find specific address
	//2. update specific peer as leader
	rep := la.parliament.FindRepresentativeByIpAddress(ipAddress)

	if rep == nil {
		return ErrNoMatchingPeerWithIpAddress
	}

	la.parliament.SetLeader(&pbft.Representative{
		ID:        rep.ID,
		IpAddress: ipAddress,
	})

	event := event.LeaderUpdated{
		LeaderId: rep.ID,
	}

	return la.eventService.Publish("leader.updated", event)

}

func (la *LeaderApi) GetParliament() *pbft.Parliament {
	return la.parliament
}

func (la *LeaderApi) GetLeader() *pbft.Leader {
	return la.parliament.Leader
}
