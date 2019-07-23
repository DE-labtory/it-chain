/*
 * Copyright 2018 DE-labtory
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
	"github.com/DE-labtory/it-chain/common"
	"github.com/DE-labtory/it-chain/common/command"
	"github.com/DE-labtory/it-chain/consensus/pbft"
	"github.com/DE-labtory/it-chain/consensus/pbft/api"
	"github.com/DE-labtory/iLogger"
	"github.com/pkg/errors"
)

var BlockSealIsNilError = errors.New("Block seal in command is nil!")

type StateStartApi struct {
	StartConsensus func(proposedBlock pbft.ProposedBlock) error
}

type StartConsensusCommandHandler struct {
	sApi *api.StateApi
}

func NewStartConsensusCommandHandler(sApi *api.StateApi) *StartConsensusCommandHandler {
	return &StartConsensusCommandHandler{
		sApi: sApi,
	}
}

func (r StartConsensusCommandHandler) HandleStartConsensusCommand(command command.StartConsensus) {

	iLogger.Info(nil, "[PBFT] Start PBFT")

	proposedBlock, err := extractProposedBlock(command)
	if err != nil {
		iLogger.Errorf(nil, "[PBFT] Extracting event is failed! - Err:[%s]", err.Error())
		return
	}

	if err = r.sApi.StartConsensus(proposedBlock); err != nil {
		iLogger.Errorf(nil, "[PBFT] Starting consensus is failed! - Err:[%s]", err.Error())
		return
	}
}

func extractProposedBlock(command command.StartConsensus) (pbft.ProposedBlock, error) {
	if command.Seal == nil {
		return pbft.ProposedBlock{}, BlockSealIsNilError
	}

	body, err := common.Serialize(command)
	if err != nil {
		return pbft.ProposedBlock{}, err
	}

	return pbft.ProposedBlock{
		Seal: command.Seal,
		Body: body,
	}, nil
}
