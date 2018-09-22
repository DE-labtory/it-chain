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
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/grpc_gateway"
	"github.com/it-chain/engine/grpc_gateway/api"
	"github.com/it-chain/iLogger"
)

type ConnectionCommandHandler struct {
	connectionApi *api.ConnectionApi
}

func NewConnectionCommandHandler(connectionApi *api.ConnectionApi) *ConnectionCommandHandler {
	return &ConnectionCommandHandler{
		connectionApi: connectionApi,
	}
}

func (d *ConnectionCommandHandler) HandleCreateConnectionCommand(createConnectionCommand command.CreateConnection) (grpc_gateway.Connection, rpc.Error) {

	connection, err := d.connectionApi.Dial(createConnectionCommand.Address)
	if err != nil {
		iLogger.Errorf(nil, "[gRPC-Gateway] Fail to dial - url [%s], Err: [%s]", createConnectionCommand.Address, err.Error())
		return grpc_gateway.Connection{}, rpc.Error{Message: err.Error()}
	}

	return connection, rpc.Error{}
}

func (d *ConnectionCommandHandler) HandleCloseConnectionCommand(closeConnectionCommand command.CloseConnection) (struct{}, rpc.Error) {

	err := d.connectionApi.CloseConnection(closeConnectionCommand.Address)
	if err != nil {
		iLogger.Errorf(nil, "[gRPC-Gateway] Fail to close - Url [%s], Err: [%s]", closeConnectionCommand.Address, err.Error())
		return struct{}{}, rpc.Error{Message: err.Error()}
	}

	return struct{}{}, rpc.Error{}
}

func (d *ConnectionCommandHandler) HandleGetConnectionListCommand(getConnectionListCommand command.GetConnectionList) (command.GetConnectionList, rpc.Error) {

	connectionList, err := d.connectionApi.GetAllConnections()
	if err != nil {
		iLogger.Errorf(nil, "[gRPC-Gateway] Fail to get connection list")
		return command.GetConnectionList{}, rpc.Error{Message: err.Error()}
	}

	return command.GetConnectionList{ConnectionList: connectionList}, rpc.Error{}
}

func (d *ConnectionCommandHandler) HandleJoinNetworkCommand(joinNetworkCommand command.JoinNetwork) (struct{}, rpc.Error) {

	err := d.connectionApi.JoinNetwork(joinNetworkCommand.Address)
	if err != nil {
		iLogger.Errorf(nil, "[gRPC-Gateway] Fail to join network - Err: [%s]", err.Error())
		return struct{}{}, rpc.Error{Message: err.Error()}
	}

	return struct{}{}, rpc.Error{}
}
