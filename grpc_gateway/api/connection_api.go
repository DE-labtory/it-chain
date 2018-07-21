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

package api

import (
	"log"

	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/engine/grpc_gateway"
)

type ConnectionApi struct {
	grpcService grpc_gateway.GrpcService
}

func NewConnectionApi(grpcService grpc_gateway.GrpcService) *ConnectionApi {
	return &ConnectionApi{
		grpcService: grpcService,
	}
}

func (c ConnectionApi) CreateConnection(address string) (grpc_gateway.Connection, error) {

	log.Printf("dialing [%s]", address)

	connection, err := c.grpcService.Dial(address)

	if err != nil {
		log.Printf("fail to dial [%s]", err)
		return grpc_gateway.Connection{}, err
	}

	return grpc_gateway.NewConnection(connection.ID, connection.Address)
}

func (c ConnectionApi) CloseConnection(connectionID string) error {

	connection := &grpc_gateway.Connection{}

	err := eventstore.Load(connection, connectionID)

	if err != nil {
		return err
	}

	c.grpcService.CloseConnection(connectionID)

	return grpc_gateway.CloseConnection(connection.ID)
}

func (c ConnectionApi) OnConnection(connection grpc_gateway.Connection) (grpc_gateway.Connection, error) {

	return grpc_gateway.NewConnection(connection.ID, connection.Address)
}

func (c ConnectionApi) OnDisconnection(connection grpc_gateway.Connection) error {

	return grpc_gateway.CloseConnection(connection.ID)
}
