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
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/engine/consensus/api"
)

type SendMsgCommandHandler struct {
	cApi api.ConsensusApi
}

func NewSendMsgCommandHandler(cApi api.ConsensusApi) SendMsgCommandHandler {
	return SendMsgCommandHandler{}
}

func (s *SendMsgCommandHandler) HandleSendPrePrepareMsg(command command.SendPrePrepareMsg) {
	representatives := make([]*consensus.Representative, 0)

	representativeList := command.RepresentativeList

	for _, r := range representativeList {
		representatives = append(representatives, &consensus.Representative{
			Id: consensus.RepresentativeId(*r),
		})
	}

	msg := consensus.PrePrepareMsg{
		ConsensusId:    consensus.NewConsensusId(command.ConsensusId),
		SenderId:       command.SenderId,
		Representative: representatives,
		ProposedBlock: consensus.ProposedBlock{
			Seal: command.Seal,
			Body: command.Body,
		},
	}

	s.cApi.ReceivePrePrepareMsg(msg)
}

func (s *SendMsgCommandHandler) HandleSendPrepareMsg(command command.SendPrepareMsg) {
	msg := consensus.PrepareMsg{
		ConsensusId: consensus.NewConsensusId(command.ConsensusId),
		SenderId:    command.SenderId,
		BlockHash:   command.BlockHash,
	}

	s.cApi.ReceivePrepareMsg(msg)
}

func (s *SendMsgCommandHandler) HandleSendCommitMsg(command command.SendCommitMsg) {
	msg := consensus.CommitMsg{
		ConsensusId: consensus.NewConsensusId(command.ConsensusId),
		SenderId:    command.SenderId,
	}

	s.cApi.ReceiveCommitMsg(msg)
}
