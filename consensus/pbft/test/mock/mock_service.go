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

type EventService struct {
	PublishFunc      func(topic string, event interface{}) error
	ConfirmBlockFunc func(block pbft.ProposedBlock) error
}

func (m EventService) Publish(topic string, event interface{}) error {
	return m.PublishFunc(topic, event)
}

func (m EventService) ConfirmBlock(block pbft.ProposedBlock) error {
	return m.ConfirmBlockFunc(block)
}

type MockConfirmService struct {
	ConfirmBlockFunc func(block pbft.ProposedBlock) error
}

func (m MockConfirmService) ConfirmBlock(block pbft.ProposedBlock) error {
	return m.ConfirmBlockFunc(block)
}

type MockPropagateService struct {
	BroadcastPrevoteMsgFunc   func(msg pbft.PrevoteMsg) error
	BroadcastProposeMsgFunc   func(msg pbft.ProposeMsg) error
	BroadcastPreCommitMsgFunc func(msg pbft.PreCommitMsg) error
}

func (m MockPropagateService) BroadcastPrevoteMsg(msg pbft.PrevoteMsg) error {
	return m.BroadcastPrevoteMsgFunc(msg)
}

func (m MockPropagateService) BroadcastProposeMsg(msg pbft.ProposeMsg) error {
	return m.BroadcastProposeMsgFunc(msg)
}
func (m MockPropagateService) BroadcastPreCommitMsg(msg pbft.PreCommitMsg) error {
	return m.BroadcastPreCommitMsgFunc(msg)
}

type MockParliamentService struct {
	RequestLeaderFunc   func() (pbft.MemberID, error)
	RequestPeerListFunc func() ([]pbft.MemberID, error)
	IsNeedConsensusFunc func() bool
}

func (m MockParliamentService) RequestLeader() (pbft.MemberID, error) {
	return m.RequestLeaderFunc()
}
func (m MockParliamentService) RequestPeerList() ([]pbft.MemberID, error) {
	return m.RequestPeerListFunc()
}
func (m MockParliamentService) IsNeedConsensus() bool {
	return m.IsNeedConsensusFunc()
}
