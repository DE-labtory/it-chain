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

	"github.com/it-chain/engine/p2p"
)

var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyConnectionId = errors.New("empty connection id proposed")
var ErrNoMatchingPeerWithIpAddress = errors.New("no matching peer with ip address")

type ILeaderApi interface {
	UpdateLeaderWithAddress(ipAddress string) error
	UpdateLeaderWithLargePeerTable(oppositePLTable p2p.PLTable) error
}

type LeaderApi struct {
	leaderService    p2p.ILeaderService
	peerQueryService p2p.PeerQueryService
}

func NewLeaderApi(leaderService p2p.ILeaderService, peerQueryService p2p.PeerQueryService) LeaderApi {

	return LeaderApi{
		leaderService:    leaderService,
		peerQueryService: peerQueryService,
	}
}

func (la *LeaderApi) UpdateLeaderWithAddress(ipAddress string) error {

	//1. loop peer list and find specific address
	//2. update specific peer as leader
	pLTable, _ := la.peerQueryService.GetPLTable()

	peers := pLTable.PeerTable

	for _, peer := range peers {

		if peer.IpAddress == ipAddress {

			la.leaderService.Set(p2p.Leader{
				LeaderId: p2p.LeaderId{
					Id: peer.PeerId.Id,
				},
			})

			return nil
		}

	}

	return ErrNoMatchingPeerWithIpAddress
}

func (la *LeaderApi) UpdateLeaderWithLargePeerTable(oppositePLTable p2p.PLTable) error {

	myPLTable, _ := la.peerQueryService.GetPLTable()

	myLeader, _ := myPLTable.GetLeader()

	if len(myPLTable.PeerTable) < len(oppositePLTable.PeerTable) {

		la.leaderService.Set(oppositePLTable.Leader)

	} else {

		la.leaderService.Set(myLeader)

	}

	return nil
}
