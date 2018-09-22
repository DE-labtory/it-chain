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

package adapter_test

import (
	"testing"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestQuerySerivce_GetRandomPeer(t *testing.T) {

	//given

	blockAdapter := mock.BlockAdapter{}

	lastBlock := blockchain.DefaultBlock{
		Seal:   []byte("seal"),
		Height: 10,
	}

	_5thBlock := blockchain.DefaultBlock{
		Seal:   []byte("seal"),
		Height: 5,
	}

	blockAdapter.GetLastBlockFromPeerFunc = func(peer blockchain.Peer) (blockchain.DefaultBlock, error) {
		return lastBlock, nil
	}

	blockAdapter.GetBlockByHeightFromPeerFunc = func(height blockchain.BlockHeight, peer blockchain.Peer) (blockchain.DefaultBlock, error) {
		return _5thBlock, nil
	}

	connectionQueryApi := mock.ConnectionQueryApi{}

	peerList := []blockchain.Peer{
		{
			PeerID:            "connection1",
			ApiGatewayAddress: "address1",
		},
		{
			PeerID:            "connection2",
			ApiGatewayAddress: "address2",
		},
		{
			PeerID:            "connection3",
			ApiGatewayAddress: "address3",
		},
	}

	connectionList := []api_gateway.Connection{
		{
			ConnectionID:      "connection1",
			ApiGatewayAddress: "address1",
		},
		{
			ConnectionID:      "connection2",
			ApiGatewayAddress: "address2",
		},
		{
			ConnectionID:      "connection3",
			ApiGatewayAddress: "address3",
		},
	}

	connectionQueryApi.GetAllConnectionListFunc = func() ([]api_gateway.Connection, error) {
		return connectionList, nil
	}

	queryService := adapter.NewQueryService(blockAdapter, connectionQueryApi)

	//when
	randomPeer, err := queryService.GetRandomPeer()

	//then
	assert.NoError(t, err)
	assert.Contains(t, peerList, randomPeer)

	//when
	retrieved_lastBlock, err := queryService.GetLastBlockFromPeer(randomPeer)

	//then
	assert.NoError(t, err)
	assert.Equal(t, lastBlock, retrieved_lastBlock)

	//when
	retrieved_5thBlock, err := queryService.GetBlockByHeightFromPeer(5, randomPeer)

	//then
	assert.NoError(t, err)
	assert.Equal(t, _5thBlock, retrieved_5thBlock)

}
