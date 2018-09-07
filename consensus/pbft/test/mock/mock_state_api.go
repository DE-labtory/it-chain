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

type StateApi struct {
	StartConsensusFunc     func(proposedBlock pbft.ProposedBlock) error
	HandleProposeMsgFunc   func(msg pbft.ProposeMsg) error
	HandlePrevoteMsgFunc   func(msg pbft.PrevoteMsg) error
	HandlePreCommitMsgFunc func(msg pbft.PreCommitMsg) error
}

func (mca *StateApi) StartConsensus(proposedBlock pbft.ProposedBlock) error {
	return mca.StartConsensusFunc(proposedBlock)
}

func (mca *StateApi) HandleProposeMsg(msg pbft.ProposeMsg) error {
	return mca.HandleProposeMsgFunc(msg)
}

func (mca *StateApi) HandlePrevoteMsg(msg pbft.PrevoteMsg) error {

	return mca.HandlePrevoteMsgFunc(msg)
}

func (mca *StateApi) HandlePreCommitMsg(msg pbft.PreCommitMsg) error {

	return mca.HandlePreCommitMsgFunc(msg)
}
