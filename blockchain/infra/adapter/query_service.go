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
	"math/rand"

	"time"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/blockchain"
)

type ConnectionQueryApi interface {
	GetAllConnectionList() ([]api_gateway.Connection, error)
	GetConnectionByID(connectionId string) (api_gateway.Connection, error)
}

type BlockAdapter interface {
	GetLastBlockFromPeer(peer blockchain.Peer) (blockchain.DefaultBlock, error)
	GetBlockByHeightFromPeer(height blockchain.BlockHeight, peer blockchain.Peer) (blockchain.DefaultBlock, error)
}

type QuerySerivce struct {
	blockAdapter       BlockAdapter
	connectionQueryApi ConnectionQueryApi
}

func NewQueryService(blockAdapter BlockAdapter, connectionQueryApi ConnectionQueryApi) *QuerySerivce {
	return &QuerySerivce{
		blockAdapter:       blockAdapter,
		connectionQueryApi: connectionQueryApi,
	}
}

func (s QuerySerivce) GetRandomPeer() (blockchain.Peer, error) {

	connectionList, err := s.connectionQueryApi.GetAllConnectionList()
	if err != nil {
		return blockchain.Peer{}, err
	}

	if len(connectionList) == 0 {
		return blockchain.Peer{}, nil
	}

	randSource := rand.NewSource(time.Now().UnixNano())
	randInstance := rand.New(randSource)
	randomIndex := randInstance.Intn(len(connectionList))
	randomPeer := toPeerFromConnection(connectionList[randomIndex])

	return randomPeer, nil
}

func (s QuerySerivce) GetLastBlockFromPeer(peer blockchain.Peer) (blockchain.DefaultBlock, error) {

	block, err := s.blockAdapter.GetLastBlockFromPeer(peer)
	if err != nil {
		return blockchain.DefaultBlock{}, err
	}

	return block, nil
}

func (s QuerySerivce) GetBlockByHeightFromPeer(height blockchain.BlockHeight, peer blockchain.Peer) (blockchain.DefaultBlock, error) {

	block, err := s.blockAdapter.GetBlockByHeightFromPeer(height, peer)
	if err != nil {
		return blockchain.DefaultBlock{}, err
	}

	return block, nil
}
