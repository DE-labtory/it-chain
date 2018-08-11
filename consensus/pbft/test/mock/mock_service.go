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

import "github.com/it-chain/engine/consensus"

type MockConfirmService struct {
	ConfirmBlockFunc func(block consensus.ProposedBlock) error
}

func (m MockConfirmService) ConfirmBlock(block consensus.ProposedBlock) error {
	return m.ConfirmBlockFunc(block)
}

type MockPropagateService struct {
	BroadcastPrepareMsgFunc    func(msg consensus.PrepareMsg) error
	BroadcastPrePrepareMsgFunc func(msg consensus.PrePrepareMsg) error
	BroadcastCommitMsgFunc     func(msg consensus.CommitMsg) error
}

func (m MockPropagateService) BroadcastPrepareMsg(msg consensus.PrepareMsg) error {
	return m.BroadcastPrepareMsgFunc(msg)
}

func (m MockPropagateService) BroadcastPrePrepareMsg(msg consensus.PrePrepareMsg) error {
	return m.BroadcastPrePrepareMsgFunc(msg)
}
func (m MockPropagateService) BroadcastCommitMsg(msg consensus.CommitMsg) error {
	return m.BroadcastCommitMsgFunc(msg)
}

type MockParliamentService struct {
	RequestLeaderFunc   func() (consensus.MemberId, error)
	RequestPeerListFunc func() ([]consensus.MemberId, error)
	IsNeedConsensusFunc func() bool
}

func (m MockParliamentService) RequestLeader() (consensus.MemberId, error) {
	return m.RequestLeaderFunc()
}
func (m MockParliamentService) RequestPeerList() ([]consensus.MemberId, error) {
	return m.RequestPeerListFunc()
}
func (m MockParliamentService) IsNeedConsensus() bool {
	return m.IsNeedConsensusFunc()
}
