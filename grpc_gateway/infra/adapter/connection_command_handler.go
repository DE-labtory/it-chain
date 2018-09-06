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
	"fmt"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/grpc_gateway"
	"github.com/it-chain/engine/grpc_gateway/api"
)

type ConnectionCommandHandler struct {
	connectionApi *api.ConnectionApi
}

func NewConnectionCommandHandler(connectionApi *api.ConnectionApi) *ConnectionCommandHandler {
	return &ConnectionCommandHandler{
		connectionApi: connectionApi,
	}
}

func (d *ConnectionCommandHandler) HandleDeployCommand(dialCommand command.Dial) (grpc_gateway.Connection, rpc.Error) {

	connection, err := d.connectionApi.CreateConnection(dialCommand.Address)
	if err != nil {
		logger.Error(nil, fmt.Sprintf("[Grpc-Gateway] fail to dial, url %s", dialCommand.Address))
		return grpc_gateway.Connection{}, rpc.Error{Message: err.Error()}
	}

	return connection, rpc.Error{}
}

func (d *ConnectionCommandHandler) HandleCloseConnectionCommand(closeConnectionCommand command.CloseConnection) (struct{}, rpc.Error) {

	err := d.connectionApi.CloseConnection(closeConnectionCommand.Address)
	if err != nil {
		logger.Error(nil, fmt.Sprintf("[Grpc-Gateway] fail to close, url %s", closeConnectionCommand.Address))
		return struct{}{}, rpc.Error{Message: err.Error()}
	}

	return struct{}{}, rpc.Error{}
}

func (d *ConnectionCommandHandler) HandleGetConnectionListCommand(getConnectionListCommand command.GetConnectionList) ([]grpc_gateway.Connection, rpc.Error) {

	connectionList, err := d.connectionApi.GetAllConnections()

	if err != nil {
		logger.Error(nil, fmt.Sprintf("[Grpc-Gateway] fail to get connection list"))
		return connectionList, rpc.Error{Message: err.Error()}
	}

	return connectionList, rpc.Error{}
}
