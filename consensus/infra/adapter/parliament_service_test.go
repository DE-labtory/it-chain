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
	"testing"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/p2p"
	"github.com/stretchr/testify/assert"
)

func TestParliamentService_RequestLeader(t *testing.T) {
	peerRepository := api_gateway.PeerRepository{}
	peerRepository.Save(p2p.Peer{
		IpAddress: "1.1.1.1",
		PeerId:    p2p.PeerId{"p1"},
	})

	ps := NewParliamentService(peerRepository)

	l, _ := ps.RequestLeader()

	assert.Equal(t, "", l)

	peerRepository.SetLeader(p2p.Peer{
		IpAddress: "2.2.2.2",
		PeerId:    p2p.PeerId{"leader"},
	})

	l, err := ps.RequestLeader()

	assert.Equal(t, "leader", l)
	assert.Nil(t, err)
}

func TestParliamentService_RequestPeerList(t *testing.T) {

}
