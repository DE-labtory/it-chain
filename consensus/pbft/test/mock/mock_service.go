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

func (m EventService) Close(){}

type ConfirmService struct {
	ConfirmBlockFunc func(block pbft.ProposedBlock) error
}

func (m ConfirmService) ConfirmBlock(block pbft.ProposedBlock) error {
	return m.ConfirmBlockFunc(block)
}

type PropagateService struct {
	BroadcastPrevoteMsgFunc   func(msg pbft.PrevoteMsg, representatives []pbft.Representative) error
	BroadcastProposeMsgFunc   func(msg pbft.ProposeMsg, representatives []pbft.Representative) error
	BroadcastPreCommitMsgFunc func(msg pbft.PreCommitMsg, representatives []pbft.Representative) error
}

func (m PropagateService) BroadcastPrevoteMsg(msg pbft.PrevoteMsg, representatives []pbft.Representative) error {
	return m.BroadcastPrevoteMsgFunc(msg, representatives)
}

func (m PropagateService) BroadcastProposeMsg(msg pbft.ProposeMsg, representatives []pbft.Representative) error {
	return m.BroadcastProposeMsgFunc(msg, representatives)
}
func (m PropagateService) BroadcastPreCommitMsg(msg pbft.PreCommitMsg, representatives []pbft.Representative) error {
	return m.BroadcastPreCommitMsgFunc(msg, representatives)
}

type ParliamentService struct {
	RequestLeaderFunc                 func() (pbft.MemberID, error)
	RequestPeerListFunc               func() ([]pbft.MemberID, error)
	IsNeedConsensusFunc               func() bool
	BuildFunc                         func() error
	SetLeaderFunc                     func(representative *pbft.Representative) error
	GetRepresentativeByIdFunc         func(id string) *pbft.Representative
	GetRepresentativeTableFunc        func() map[string]*pbft.Representative
	GetParliamentFunc                 func() *pbft.Parliament
	GetLeaderFunc                     func() *pbft.Leader
	FindRepresentativeByIpAddressFunc func(ipAddress string) *pbft.Representative
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
func (m ParliamentService) Build() error {
	return m.BuildFunc()
}
func (m ParliamentService) SetLeader(representative *pbft.Representative) error {
	return m.SetLeaderFunc(representative)
}
func (m ParliamentService) GetRepresentativeById(id string) *pbft.Representative {
	return m.GetRepresentativeByIdFunc(id)
}
func (m ParliamentService) GetRepresentativeTable() map[string]*pbft.Representative {
	return m.GetRepresentativeTableFunc()
}
func (m ParliamentService) GetParliament() *pbft.Parliament {
	return m.GetParliamentFunc()
}
func (m ParliamentService) GetLeader() *pbft.Leader {
	return m.GetLeaderFunc()
}
func (m ParliamentService) FindRepresentativeByIpAddress(ipAddress string) *pbft.Representative {
	return m.FindRepresentativeByIpAddressFunc(ipAddress)
}
