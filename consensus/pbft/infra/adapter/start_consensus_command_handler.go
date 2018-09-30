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
	"github.com/it-chain/iLogger"
	"github.com/pkg/errors"
)

var BlockSealIsNilError = errors.New("Block seal in command is nil!")

type StateStartApi struct {
	StartConsensus func(proposedBlock pbft.ProposedBlock) error
}

type StartConsensusCommandHandler struct {
	sApi StateStartApi
}

func NewStartConsensusCommandHandler(sApi StateStartApi) *StartConsensusCommandHandler {
	return &StartConsensusCommandHandler{
		sApi: sApi,
	}
}

func (r StartConsensusCommandHandler) HandleStartConsensusCommand(created event.BlockCreated) {

	seal := created.Seal
	txList := created.TxList

	proposedBlock, err := extractProposedBlock(seal, txList)
	if err != nil {
		iLogger.Errorf(nil, "[PBFT] Extracting event is failed! - %s", err.Error())
		return
	}

	if err = r.sApi.StartConsensus(proposedBlock); err != nil {
		iLogger.Errorf(nil, "[PBFT] Starting consensus is failed! - %s", err.Error())
		return
	}
}

func extractProposedBlock(Seal []byte, TxList []event.Tx) (pbft.ProposedBlock, error) {
	if Seal == nil {
		return pbft.ProposedBlock{}, BlockSealIsNilError
	}

	body, err := common.Serialize(TxList)
	if err != nil {
		return pbft.ProposedBlock{}, err
	}

	return pbft.ProposedBlock{
		Seal: Seal,
		Body: body,
	}, nil
}
