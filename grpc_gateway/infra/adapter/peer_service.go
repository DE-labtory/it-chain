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
	"encoding/json"

	"github.com/it-chain/engine/grpc_gateway"
	"gopkg.in/resty.v1"
)

type PeerAdapter interface {
	GetAllPeerList(address string) ([]grpc_gateway.Connection, error)
}

type HttpPeerAdapter struct{}

func NewHttpPeerAdapter() *HttpPeerAdapter {
	return &HttpPeerAdapter{}
}

func (HttpPeerAdapter) GetAllPeerList(address string) ([]grpc_gateway.Connection, error) {
	resp, err := resty.R().Get("http://" + address + "/connections")

	if err != nil {
		return []grpc_gateway.Connection{}, err
	}

	connectionList := []grpc_gateway.Connection{}
	if err := json.Unmarshal(resp.Body(), &connectionList); err != nil {
		return []grpc_gateway.Connection{}, err
	}

	return connectionList, nil
}

type PeerService struct {
	PeerAdapter
}

func NewPeerService(peerAdapter PeerAdapter) *PeerService {
	return &PeerService{
		PeerAdapter: peerAdapter,
	}
}

func (p PeerService) GetAllPeerList(connection grpc_gateway.Connection) ([]grpc_gateway.Connection, error) {
	return p.PeerAdapter.GetAllPeerList(connection.ApiGatewayAddress)
}
