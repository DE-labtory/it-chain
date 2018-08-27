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
	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/api"
)

var ErrLeaderInfoDeliver = errors.New("leader info deliver failed")
var ErrPeerListDeliver = errors.New("peer list deliver failed")
var ErrPeerDeliver = errors.New("peer deliver failed")
var ErrUnmarshal = errors.New("error during unmarshal")

type GrpcCommandHandler struct {
	leaderApi        api.ILeaderApi
	electionService  *p2p.ElectionService
	communicationApi CommunicationApi // api.CommunicationApi
	pLTableService   p2p.PLTableService
}

func NewGrpcCommandHandler(
	leaderApi api.ILeaderApi,
	electionService *p2p.ElectionService, communicationApi CommunicationApi,
	pLTableService p2p.PLTableService) GrpcCommandHandler {
	return GrpcCommandHandler{
		leaderApi:        leaderApi,
		electionService:  electionService,
		communicationApi: communicationApi,
		pLTableService:   pLTableService,
	}
}

func (gch *GrpcCommandHandler) HandleMessageReceive(command command.ReceiveGrpc) error {

	switch command.Protocol {

	case "PLTableDeliverProtocol": //receive peer table

		//1. receive peer table
		pLTable, _ := gch.pLTableService.GetPLTableFromCommand(command)

		//2. update leader and peer list by info of node which has longer peer list
		gch.leaderApi.UpdateLeaderWithLargePeerTable(pLTable)

		//3. dial according to peer table
		gch.communicationApi.DialToUnConnectedNode(pLTable.PeerTable)

		break

	case "RequestVoteProtocol":
		logger.Infof(nil, "handling request vote from process: %v", gch.electionService.Election.GetIpAddress())
		gch.electionService.Vote(command.ConnectionID)

	case "VoteLeaderProtocol":
		//	1. if candidate, reset left time
		//	2. count up
		//	3. if counted is same with num of peer-1 set leader and publish

		logger.Infof(nil, "received VoteLeaderProtocol command:", command)
		logger.Infof(nil, "received VoteLeaderProtocol current election :", gch.electionService.Election)
		gch.electionService.DecideToBeLeader(command)

	case "UpdateLeaderProtocol":

		// if received leader is not what i voted for, return nil
		if gch.electionService.Election.GetCandidate().PeerId.Id != command.ConnectionID {
			return nil
		}

		toBeLeader := &p2p.UpdateLeaderMessage{}
		err := common.Deserialize(command.Body, toBeLeader)

		logger.Infof(nil, "update leader with", toBeLeader.Peer)
		if err != nil {
			return err
		}

		gch.leaderApi.UpdateLeaderWithAddress(toBeLeader.Peer.IpAddress)
	}

	return nil
}

//
