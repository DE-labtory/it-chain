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

package p2p

import (
	"encoding/json"
	"errors"

	"github.com/it-chain/engine/common/command"
)

var ErrEmptyLeaderId = errors.New("empty leader id")
var ErrEmptyPeerTable = errors.New("empty peer list")

type PLTable struct {
	Leader    Leader
	PeerTable map[string]Peer
}

func NewPLTable(leader Leader, peerTable map[string]Peer) *PLTable {

	return &PLTable{
		Leader:    leader,
		PeerTable: peerTable,
	}
}

func (pt *PLTable) GetLeader() (Leader, error) {

	if pt.Leader.LeaderId.Id == "" {

		return pt.Leader, ErrEmptyLeaderId
	}

	return pt.Leader, nil
}

func (pt *PLTable) GetPeerTable() (map[string]Peer, error) {

	if len(pt.PeerTable) == 0 {

		return pt.PeerTable, ErrEmptyPeerTable
	}

	return pt.PeerTable, nil
}

type PLTableServiceImpl struct{}

func (plts *PLTableServiceImpl) GetPLTableFromCommand(command command.ReceiveGrpc) (PLTable, error) {

	peerTable := PLTable{}

	if err := json.Unmarshal(command.Body, &peerTable); err != nil {
		//todo error 처리
		return PLTable{}, nil
	}

	return peerTable, nil
}
