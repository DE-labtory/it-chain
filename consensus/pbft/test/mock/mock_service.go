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
	PublishFunc func(topic string, event interface{}) error
}

func (m EventService) Publish(topic string, event interface{}) error {
	return m.PublishFunc(topic, event)
}

type ConfirmService struct {
	ConfirmBlockFunc func(block pbft.ProposedBlock) error
}

func (m ConfirmService) ConfirmBlock(block pbft.ProposedBlock) error {
	return m.ConfirmBlockFunc(block)
}

type PropagateService struct {
	BroadcastPrepareMsgFunc    func(msg pbft.PrepareMsg) error
	BroadcastPrePrepareMsgFunc func(msg pbft.PrePrepareMsg) error
	BroadcastCommitMsgFunc     func(msg pbft.CommitMsg) error
}

func (m PropagateService) BroadcastPrepareMsg(msg pbft.PrepareMsg) error {
	return m.BroadcastPrepareMsgFunc(msg)
}

func (m PropagateService) BroadcastPrePrepareMsg(msg pbft.PrePrepareMsg) error {
	return m.BroadcastPrePrepareMsgFunc(msg)
}
func (m PropagateService) BroadcastCommitMsg(msg pbft.CommitMsg) error {
	return m.BroadcastCommitMsgFunc(msg)
}

type ParliamentService struct {
	RequestLeaderFunc   func() (pbft.MemberID, error)
	RequestPeerListFunc func() ([]pbft.MemberID, error)
	IsNeedConsensusFunc func() bool
}

func (m ParliamentService) RequestLeader() (pbft.MemberID, error) {
	return m.RequestLeaderFunc()
}
func (m ParliamentService) RequestPeerList() ([]pbft.MemberID, error) {
	return m.RequestPeerListFunc()
}
func (m ParliamentService) IsNeedConsensus() bool {
	return m.IsNeedConsensusFunc()
}
