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

package api_gateway_test

import (
	"testing"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/common/event"
	"github.com/stretchr/testify/assert"
)

func TestPeerRepository_SetInvalidLeader(t *testing.T) {
	//given
	peerRepository := api_gateway.NewPeerRepository()

	//when
	err := peerRepository.SetLeader("1")

	//then
	assert.Equal(t, err, api_gateway.ErrPeerDoesNotExists)
}

func TestPeerRepository_SetLeader_GetLeader(t *testing.T) {

	//given
	peerRepository := api_gateway.NewPeerRepository()
	peerRepository.Save(api_gateway.Peer{
		ID:   "123",
		Role: api_gateway.Member,
	})

	//when
	err := peerRepository.SetLeader("123")
	assert.NoError(t, err)

	//then
	leader, err := peerRepository.GetLeader()
	assert.NoError(t, err)
	assert.Equal(t, leader.ID, "123")
}

func TestPeerRepository_FindAll(t *testing.T) {

	//given
	peerRepository := api_gateway.NewPeerRepository()
	peer1 := api_gateway.Peer{
		ID:   "123",
		Role: api_gateway.Member,
	}
	peer2 := api_gateway.Peer{
		ID:   "124",
		Role: api_gateway.Member,
	}

	peerRepository.Save(peer1)
	peerRepository.Save(peer2)

	//when, then
	assert.Contains(t, peerRepository.FindAll(), peer1)
	assert.Contains(t, peerRepository.FindAll(), peer2)
}

func TestPeerRepository_FindById(t *testing.T) {
	//given
	peerRepository := api_gateway.NewPeerRepository()
	peer1 := api_gateway.Peer{
		ID:   "123",
		Role: api_gateway.Member,
	}
	peerRepository.Save(peer1)

	//when
	peer, err := peerRepository.FindById("123")

	//then
	assert.NoError(t, err)
	assert.Equal(t, peer, peer1)
}

func TestPeerRepository_Remove(t *testing.T) {
	//given
	peerRepository := api_gateway.NewPeerRepository()
	peer1 := api_gateway.Peer{
		ID:   "123",
		Role: api_gateway.Member,
	}
	peerRepository.Save(peer1)

	//when
	peerRepository.Remove("123")

	//then
	_, err := peerRepository.FindById("123")
	assert.Equal(t, err, api_gateway.ErrPeerDoesNotExists)
}

func TestPeerRepository_Save(t *testing.T) {
	//given
	peerRepository := api_gateway.NewPeerRepository()
	peer1 := api_gateway.Peer{
		ID:   "123",
		Role: api_gateway.Member,
	}
	peerRepository.Save(peer1)

	//when
	peer, err := peerRepository.FindById("123")

	//then
	assert.NoError(t, err)
	assert.Equal(t, peer, peer1)
}

func TestConnectionEventListener_HandleConnectionCreatedEvent(t *testing.T) {

	//given
	peerRepository := api_gateway.NewPeerRepository()
	listener := api_gateway.NewConnectionEventListener(peerRepository)
	listener.HandleConnectionCreatedEvent(event.ConnectionCreated{
		ConnectionID:       "0",
		GrpcGatewayAddress: "address",
		ApiGatewayAddress:  "123",
	})

	//when
	peer := peerRepository.FindAll()[0]

	//then
	assert.Equal(t, peer.Role, api_gateway.Member)
	assert.Equal(t, peer.ID, "0")
	assert.Equal(t, peer.GrpcGatewayAddress, "address")
	assert.Equal(t, peer.ApiGatewayAddress, "123")
}

func TestConnectionEventListener_HandleConnectionClosedEvent(t *testing.T) {

	//given
	peerRepository := api_gateway.NewPeerRepository()
	peerRepository.Save(api_gateway.Peer{
		ID: "123",
	})
	listener := api_gateway.NewConnectionEventListener(peerRepository)

	//when
	listener.HandleConnectionClosedEvent(event.ConnectionClosed{
		ConnectionID: "123",
	})

	//then
	_, err := peerRepository.FindById("123")
	assert.Equal(t, err, api_gateway.ErrPeerDoesNotExists)
}

func TestLeaderUpdateEventListener_HandleLeaderUpdatedEvent(t *testing.T) {
	peerRepository := api_gateway.NewPeerRepository()
	peerRepository.Save(api_gateway.Peer{
		ID:   "123",
		Role: api_gateway.Member,
	})

	listener := api_gateway.NewLeaderUpdateEventListener(peerRepository)
	listener.HandleLeaderUpdatedEvent(event.LeaderUpdated{
		LeaderId: "123",
	})

	peer, err := peerRepository.GetLeader()
	assert.NoError(t, err)

	assert.Equal(t, peer.ID, "123")
}

func TestPeerQueryApi_GetPeerByID(t *testing.T) {
	peerRepository := api_gateway.NewPeerRepository()
	peerRepository.Save(api_gateway.Peer{
		ID:   "123",
		Role: api_gateway.Member,
	})

	peerQueryApi := api_gateway.NewPeerQueryApi(peerRepository)

	peer, err := peerQueryApi.GetPeerByID("123")
	assert.NoError(t, err)
	assert.Equal(t, peer.ID, "123")
	assert.Equal(t, peer.Role, api_gateway.Member)
}

func TestPeerQueryApi_GetAllPeerList(t *testing.T) {

	//given
	peerRepository := api_gateway.NewPeerRepository()
	peerRepository.Save(api_gateway.Peer{
		ID:   "123",
		Role: api_gateway.Member,
	})

	peerRepository.Save(api_gateway.Peer{
		ID:   "1234",
		Role: api_gateway.Member,
	})
	peerQueryApi := api_gateway.NewPeerQueryApi(peerRepository)

	//when
	peerList := peerQueryApi.GetAllPeerList()

	//then
	assert.Contains(t, peerList, api_gateway.Peer{
		ID:   "123",
		Role: api_gateway.Member,
	})
	assert.Contains(t, peerList, api_gateway.Peer{
		ID:   "1234",
		Role: api_gateway.Member,
	})
}

func TestPeerQueryApi_GetLeader(t *testing.T) {

	//given
	peerRepository := api_gateway.NewPeerRepository()
	peerRepository.Save(api_gateway.Peer{
		ID:   "123",
		Role: api_gateway.Member,
	})
	peerRepository.Save(api_gateway.Peer{
		ID:   "1234",
		Role: api_gateway.Leader,
	})
	peerQueryApi := api_gateway.NewPeerQueryApi(peerRepository)

	//when
	leader, err := peerQueryApi.GetLeader()

	//then
	assert.NoError(t, err)
	assert.Equal(t, leader.ID, "1234")
}

func TestPeerQueryApi_GetNetwork(t *testing.T) {

	//given
	peerRepository := api_gateway.NewPeerRepository()
	peerRepository.Save(api_gateway.Peer{
		ID:   "123",
		Role: api_gateway.Member,
	})
	peerRepository.Save(api_gateway.Peer{
		ID:   "1234",
		Role: api_gateway.Leader,
	})
	peerQueryApi := api_gateway.NewPeerQueryApi(peerRepository)

	//when
	network, err := peerQueryApi.GetNetwork()

	//then
	assert.NoError(t, err)
	assert.Equal(t, network.Leader, api_gateway.Peer{
		ID:   "1234",
		Role: api_gateway.Leader,
	})

	assert.Contains(t, network.Members, api_gateway.Peer{
		ID:   "1234",
		Role: api_gateway.Leader,
	})
}
