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
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/api"
	"github.com/pkg/errors"
)

var BlockSealIsNilError = errors.New("Block seal in command is nil!")

type StartConsensusCommandHandler struct {
	sApi api.StateApi
}

func NewStartConsensusCommandHandler(sApi api.StateApi) *StartConsensusCommandHandler {
	return &StartConsensusCommandHandler{
		sApi: sApi,
	}
}

func (r StartConsensusCommandHandler) HandleStartConsensusCommand(startConsensusCommand command.StartConsensus) (bool, rpc.Error) {
	seal := startConsensusCommand.Seal
	txList := startConsensusCommand.TxList

	proposedBlock, err := extractProposedBlock(seal, txList)
	if err != nil {
		return false, rpc.Error{Message: err.Error()}
	}

	if err = r.sApi.StartConsensus(proposedBlock); err != nil {
		return false, rpc.Error{Message: err.Error()}
	}

	return true, rpc.Error{}
}

func extractProposedBlock(Seal []byte, TxList []command.Tx) (pbft.ProposedBlock, error) {
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
