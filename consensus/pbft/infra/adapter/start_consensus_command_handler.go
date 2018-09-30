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
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/api"
	"github.com/it-chain/iLogger"
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

func (r StartConsensusCommandHandler) HandleStartConsensusCommand(created event.BlockCreated) {

	iLogger.Info(nil, "[PBFT] Start PBFT")

	proposedBlock, err := extractProposedBlock(created)
	if err != nil {
		iLogger.Errorf(nil, "[PBFT] Extracting event is failed! - %s", err.Error())
		return
	}

	if err = r.sApi.StartConsensus(proposedBlock); err != nil {
		iLogger.Errorf(nil, "[PBFT] Starting consensus is failed! - %s", err.Error())
		return
	}
}

func extractProposedBlock(created event.BlockCreated) (pbft.ProposedBlock, error) {
	if created.Seal == nil {
		return pbft.ProposedBlock{}, BlockSealIsNilError
	}

	body, err := common.Serialize(created)
	if err != nil {
		return pbft.ProposedBlock{}, err
	}

	return pbft.ProposedBlock{
		Seal: created.Seal,
		Body: body,
	}, nil
}
