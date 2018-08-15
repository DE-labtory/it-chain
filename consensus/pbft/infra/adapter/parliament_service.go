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

package adapter

import (
	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/consensus/pbft"
)

type ParliamentService struct {
	pQuery api_gateway.PeerQueryApi
}

func NewParliamentService(api api_gateway.PeerQueryApi) *ParliamentService {
	return &ParliamentService{
		pQuery: api,
	}
}

func (ps *ParliamentService) RequestLeader() (pbft.MemberId, error) {
	l, err := ps.pQuery.GetLeader()

	if err != nil {
		return "", err
	}

	return pbft.MemberId(l.GetID()), nil
}

func (ps *ParliamentService) RequestPeerList() ([]pbft.MemberId, error) {
	pl, err := ps.pQuery.GetPeerList()

	if err != nil {
		return nil, err
	}

	peerList := make([]pbft.MemberId, 0)

	for _, p := range pl {
		peerList = append(peerList, pbft.MemberId(p.GetID()))
	}

	return peerList, nil
}

func (p *ParliamentService) IsNeedConsensus() bool {
	peerList, err := p.RequestPeerList()

	if err != nil {
		return false
	}

	numOfMember := 0
	numOfMember = numOfMember + len(peerList)

	if numOfMember >= 4 {
		return true
	}

	return false
}
