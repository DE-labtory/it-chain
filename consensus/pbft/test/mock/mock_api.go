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

package mock

import (
	"github.com/it-chain/engine/consensus"
)

type MockConsensusApi struct {
	StartConsensusFunc       func(userId consensus.MemberId, proposedBlock consensus.ProposedBlock) error
	ReceivePrePrepareMsgFunc func(msg consensus.PrePrepareMsg) error
	ReceivePrepareMsgFunc    func(msg consensus.PrepareMsg) error
	ReceiveCommitMsgFunc     func(msg consensus.CommitMsg) error
}

func (mca *MockConsensusApi) StartConsensus(userId consensus.MemberId, proposedBlock consensus.ProposedBlock) error {
	return mca.StartConsensus(userId, proposedBlock)
}

func (mca *MockConsensusApi) ReceivePrePrepareMsg(msg consensus.PrePrepareMsg) error {
	return mca.ReceivePrePrepareMsgFunc(msg)
}

func (mca *MockConsensusApi) ReceivePrepareMsg(msg consensus.PrepareMsg) error {

	return mca.ReceivePrepareMsgFunc(msg)
}

func (mca *MockConsensusApi) ReceiveCommitMsg(msg consensus.CommitMsg) error {

	return mca.ReceiveCommitMsgFunc(msg)
}
