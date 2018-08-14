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
	"github.com/it-chain/engine/consensus/pbft"
)

type MockStateApi struct {
	StartConsensusFunc      func(userId pbft.MemberID, proposedBlock pbft.ProposedBlock) error
	HandlePrePrepareMsgFunc func(msg pbft.PrePrepareMsg) error
	HandlePrepareMsgFunc    func(msg pbft.PrepareMsg) error
	HandleCommitMsgFunc     func(msg pbft.CommitMsg) error
}

func (mca *MockStateApi) StartConsensus(userId pbft.MemberID, proposedBlock pbft.ProposedBlock) error {
	return mca.StartConsensus(userId, proposedBlock)
}

func (mca *MockStateApi) HandlePrePrepareMsg(msg pbft.PrePrepareMsg) error {
	return mca.HandlePrePrepareMsgFunc(msg)
}

func (mca *MockStateApi) HandlePrepareMsg(msg pbft.PrepareMsg) error {

	return mca.HandlePrepareMsgFunc(msg)
}

func (mca *MockStateApi) HandleCommitMsg(msg pbft.CommitMsg) error {

	return mca.HandleCommitMsgFunc(msg)
}
