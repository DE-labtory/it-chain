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
	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/midgard"
)

type ConsensusApi struct {
	eventRepository *midgard.Repository
}

// todo : Event Sourcing 첨가

func (cApi ConsensusApi) StartConsensus(userId consensus.MemberId, block consensus.ProposedBlock) error {
	return nil
}

func (cApi ConsensusApi) ReceivePrePrepareMsg(msg consensus.PrePrepareMsg) {
	return
}

func (cApi ConsensusApi) ReceivePrepareMsg(msg consensus.PrepareMsg) error{
	return nil
}

func (cApi ConsensusApi) ReceiveCommitMsg(msg consensus.CommitMsg) {
	return
}
