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

import "github.com/it-chain/engine/p2p"

type MockPLTableApi struct {
	getPLTableFunc func() p2p.PLTable
}

func (mplta *MockPLTableApi) GetPLTable() p2p.PLTable {

	return mplta.getPLTableFunc()
}

type MockPeerApi struct {
	SaveFunc   func(peer p2p.Peer) error
	RemoveFunc func(peerId p2p.PeerId) error
}

func (mpa *MockPeerApi) Save(peer p2p.Peer) error {
	return mpa.SaveFunc(peer)
}

func (mpa *MockPeerApi) Remove(peerId p2p.PeerId) error {
	return mpa.RemoveFunc(peerId)
}

type MockLeaderApi struct {
}

func (mla *MockLeaderApi) UpdateLeaderWithAddress(ipAddress string) error {
	return mla.UpdateLeaderWithAddress(ipAddress)
}

func (mla *MockLeaderApi) UpdateLeaderWithLargePeerTable(oppositePLTable p2p.PLTable) error {
	return mla.UpdateLeaderWithLargePeerTable(oppositePLTable)
}

type MockCommunicationApi struct {
	DeliverPLTableFunc       func(connectionId string) error
	DialToUnConnectedNodeFuc func(peerTable map[string]p2p.Peer) error
}

func (mca *MockCommunicationApi) DeliverPLTable(connectionId string) error {

	return mca.DeliverPLTableFunc(connectionId)
}

func (mca *MockCommunicationApi) DialToUnConnectedNode(peerTable map[string]p2p.Peer) error {

	return mca.DialToUnConnectedNodeFuc(peerTable)
}
