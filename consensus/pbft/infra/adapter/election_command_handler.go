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
	"errors"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/api"
)

var ErrLeaderInfoDeliver = errors.New("leader info deliver failed")
var ErrPeerListDeliver = errors.New("peer list deliver failed")
var ErrPeerDeliver = errors.New("peer deliver failed")
var ErrUnmarshal = errors.New("error during unmarshal")

type ElectionCommandHandler struct {
	leaderApi   *api.LeaderApi
	electionApi api.ElectionApi
}

func NewElectionCommandHandler(
	leaderApi *api.LeaderApi,
	electionApi *api.ElectionApi) *ElectionCommandHandler {
	return &ElectionCommandHandler{
		leaderApi:   leaderApi,
		electionApi: *electionApi,
	}
}

func (gch *ElectionCommandHandler) HandleMessageReceive(command command.ReceiveGrpc) error {

	switch command.Protocol {

	case "RequestVoteProtocol":
		logger.Infof(nil, "[consensus] handling request vote from process: %v", gch.electionApi.GetIpAddress())

		err := gch.electionApi.Vote(command.ConnectionID)

		if err != nil {
			return err
		}

	case "VoteLeaderProtocol":
		logger.Infof(nil, "[consensus] voted from process: %v", command.ConnectionID)

		//	1. if candidate, reset left time
		//	2. count up
		//	3. if counted is same with num of peer-1 set leader and publish

		logger.Infof(nil, "[consensus] received VoteLeaderProtocol command:", command)

		err := gch.electionApi.DecideToBeLeader()

		if err != nil {
			return err
		}

	case "UpdateLeaderProtocol":
		// if received leader is not what i voted for, return nil
		if gch.electionApi.GetCandidate().ID != command.ConnectionID {
			return nil
		}

		toBeLeader := &pbft.UpdateLeaderMessage{}
		err := common.Deserialize(command.Body, toBeLeader)

		logger.Infof(nil, "[consensus] update leader with", toBeLeader.Representative)
		if err != nil {
			return err
		}

		err2 := gch.leaderApi.UpdateLeaderWithAddress(toBeLeader.Representative.IpAddress)

		if err2 != nil {
			return err2
		}
	}

	return nil
}

//
