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
)

type ParliamentService struct {
	pQuery api_gateway.PeerQueryApi
}

func NewParliamentService(peerRepository *api_gateway.PeerRepository) *ParliamentService {
	return &ParliamentService{
		pQuery: api_gateway.NewPeerQueryApi(peerRepository),
	}
}

func (ps *ParliamentService) RequestLeader() (string, error) {
	l, err := ps.pQuery.GetLeader()

	if err != nil {
		return "", err
	}

	return l.GetID(), nil
}

func (ps *ParliamentService) RequestPeerList() ([]string, error) {
	pl, err := ps.pQuery.GetPeerList()

	if err != nil {
		return nil, err
	}

	strPL := make([]string, 0)

	for _, p := range pl {
		strPL = append(strPL, p.GetID())
	}

	return strPL, nil
}
