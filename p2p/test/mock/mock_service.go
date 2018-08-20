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
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/p2p"
)

type MockPeerService struct {
	SaveFunc   func(peer p2p.Peer) error
	RemoveFunc func(peerId p2p.PeerId) error
}

func (m *MockPeerService) Save(peer p2p.Peer) error {
	return m.SaveFunc(peer)
}

func (m *MockPeerService) Remove(peerId p2p.PeerId) error {
	return m.RemoveFunc(peerId)
}

type MockLeaderService struct {
	SetFunc func(leader p2p.Leader) error
}

func (ls *MockLeaderService) Set(leader p2p.Leader) error {

	return ls.SetFunc(leader)
}

type MockCommunicationService struct {
	DialFunc           func(ipAddress string) error
	DeliverPLTableFunc func(connectionId string, pLTable p2p.PLTable) error
}

func (mcs *MockCommunicationService) Dial(ipAddress string) error {

	return mcs.DialFunc(ipAddress)
}

func (mcs *MockCommunicationService) DeliverPLTable(connectionId string, pLTable p2p.PLTable) error {

	return mcs.DeliverPLTableFunc(connectionId, pLTable)
}

type MockPeerQueryService struct {
	GetPLTableFunc        func() (p2p.PLTable, error)
	GetPeerListFunc       func() ([]p2p.Peer, error)
	GetLeaderFunc         func() (p2p.Leader, error)
	FindPeerByIdFunc      func(peerId p2p.PeerId) (p2p.Peer, error)
	FindPeerByAddressFunc func(ipAddress string) (p2p.Peer, error)
}

func (mpltqs *MockPeerQueryService) GetPLTable() (p2p.PLTable, error) {

	return mpltqs.GetPLTableFunc()
}

func (mpltqs *MockPeerQueryService) GetPeerList() ([]p2p.Peer, error) {

	return mpltqs.GetPeerListFunc()
}

func (mpltqs *MockPeerQueryService) GetLeader() (p2p.Leader, error) {

	return mpltqs.GetLeaderFunc()
}

func (mpltqs *MockPeerQueryService) FindPeerById(peerId p2p.PeerId) (p2p.Peer, error) {

	return mpltqs.FindPeerByIdFunc(peerId)
}

func (mpltqs *MockPeerQueryService) FindPeerByAddress(ipAddress string) (p2p.Peer, error) {

	return mpltqs.FindPeerByAddressFunc(ipAddress)
}

type MockPLTableService struct {
	GetPLTableFromCommandFunc func(command command.ReceiveGrpc) (p2p.PLTable, error)
}

func (mplts *MockPLTableService) GetPLTableFromCommand(command command.ReceiveGrpc) (p2p.PLTable, error) {

	return mplts.GetPLTableFromCommandFunc(command)
}

type MockElectionService struct {
	VoteFunc             func(connectionId string) error
	BroadcastLeaderFunc  func(peer p2p.Peer) error
	DecideToBeLeaderFunc func(command command.ReceiveGrpc) error
	ElectLeaderWithRaft  func() error
}

func (mes *MockElectionService) Vote(connectionId string) error {

	return mes.VoteFunc(connectionId)

}

func (mes *MockElectionService) BroadcastLeader(peer p2p.Peer) error {

	return mes.BroadcastLeaderFunc(peer)

}

func (mes *MockElectionService) DecideToBeLeader(command command.ReceiveGrpc) error {

	return mes.DecideToBeLeaderFunc(command)

}

func (mes *MockElectionService) ElectLEaderWithRaft() error {

	return mes.ElectLeaderWithRaft()
}

type MockPublishService struct {
	PeerCreatedFunc   func(peer p2p.Peer) error
	PeerDeletedFunc   func(peerId p2p.PeerId) error
	LeaderUpdatedFunc func(leader p2p.Leader) error
}

func (mp *MockPublishService) PeerCreated(peer p2p.Peer) error {
	return mp.PeerCreatedFunc(peer)
}

func (mp *MockPublishService) PeerDeleted(peerId p2p.PeerId) error {
	return mp.PeerDeletedFunc(peerId)
}

func (mp *MockPublishService) LeaderUpdated(leader p2p.Leader) error {
	return mp.LeaderUpdatedFunc(leader)
}

type MockEventService struct {
	PublishFunc func(topic string, event interface{}) error
}

func (m *MockEventService) Publish(topic string, event interface{}) error {
	return m.PublishFunc(topic, event)
}
