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
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/api"
	"github.com/it-chain/iLogger"
)

var ErrLeaderInfoDeliver = errors.New("leader info deliver failed")
var ErrPeerListDeliver = errors.New("peer list deliver failed")
var ErrPeerDeliver = errors.New("peer deliver failed")
var ErrUnmarshal = errors.New("error during unmarshal")

type ElectionCommandHandler struct {
	parliamentApi *api.ParliamentApi
	electionApi   *api.ElectionApi
}

func NewElectionCommandHandler(parliamentApi *api.ParliamentApi, electionApi *api.ElectionApi) *ElectionCommandHandler {
	return &ElectionCommandHandler{
		parliamentApi: parliamentApi,
		electionApi:   electionApi,
	}
}

func (e *ElectionCommandHandler) HandleMessageReceive(command command.ReceiveGrpc) error {

	switch command.Protocol {

	case "RequestVoteProtocol":
		message := &pbft.RequestVoteMessage{}
		deserializeErr := common.Deserialize(command.Body, message)
		if deserializeErr != nil {
			return deserializeErr
		}

		if e.electionApi.ElectionService.Voted {
			iLogger.Info(nil, "already voted!")
			return nil
		}

		err := e.electionApi.Vote(command.ConnectionID)

		if err != nil {
			return err
		}

	case "VoteLeaderProtocol":
		iLogger.Infof(nil, "[PBFT] Receive VoteLeaderProtocol")
		if err := e.electionApi.DecideToBeLeader(); err != nil {
			iLogger.Errorf(nil, "Err %s", err.Error())
		}

	case "UpdateLeaderProtocol":
		if e.electionApi.GetCandidate().ID != command.ConnectionID {
			return nil
		}

		e.electionApi.EndRaft()
		toBeLeader := &pbft.UpdateLeaderMessage{}
		if err := common.Deserialize(command.Body, toBeLeader); err != nil {
			iLogger.Errorf(nil, "Err %s", err.Error())
		}

		if err := e.parliamentApi.UpdateLeader(toBeLeader.Representative.ID); err != nil {
			iLogger.Errorf(nil, "Err %s", err.Error())
		}
	}

	return nil
}
