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

import "github.com/it-chain/engine/consensus"

// todo : blockchain으로 부터 받은 command를 처리하는 service
type RequestedConsensusCommandHandler struct {
	consensus         *consensus.Consensus
	parliamentService ParliamentService
}

func NewRequestedConsensusSCommandHandler(c *consensus.Consensus, service ParliamentService) RequestedConsensusCommandHandler {
	return RequestedConsensusCommandHandler{
		consensus:         c,
		parliamentService: service,
	}
}

// todo : requested consensus command 가 아직 정의 되지 않음
func (r RequestedConsensusCommandHandler) HandleConsensusRequestCommand() error {
	parliament, err := r.parliamentService.RequestPeerList()

	if err != nil {
		return err
	}

	// todo : command에서 proposed block을 복구해야함
	proposedBlock := consensus.ProposedBlock{}

	c, err := consensus.CreateConsensus(parliament, proposedBlock)

	if err != nil {
		return err
	}

	// c라는 consensus 객체를 주입
	r.consensus = c
	return nil
}
