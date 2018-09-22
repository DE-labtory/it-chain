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
	"strings"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/grpc_gateway"
	"github.com/it-chain/engine/grpc_gateway/api"
	"github.com/it-chain/iLogger"
)

type GrpcMessageHandler struct {
	connectionApi *api.ConnectionApi
	messageApi    *api.MessageApi
}

func NewGrpcMessageHandler(connectionApi *api.ConnectionApi, messageApi *api.MessageApi) *GrpcMessageHandler {
	return &GrpcMessageHandler{
		connectionApi: connectionApi,
		messageApi:    messageApi,
	}
}

func (g GrpcMessageHandler) HandleMessageReceiveCommand(command command.ReceiveGrpc) {

	protocol := command.Protocol
	body := command.Body

	iLogger.Infof(nil, "[gRPC-Gateway] Received gRPC message - Protocol [%s]", protocol)

	switch protocol {
	case grpc_gateway.RequestPeerList:
		g.connectionApi.HandleRequestPeerList(command.ConnectionID)

	case grpc_gateway.ResponsePeerList:
		connectionList := []grpc_gateway.Connection{}

		if err := common.Deserialize(body, &connectionList); err != nil {
			iLogger.Errorf(nil, "[gRPC-Gateway] Fail to deserialize grpcMessage - Err: [%s]", err.Error())
			return
		}

		g.connectionApi.DialConnectionList(connectionList)
	}
}

func (g GrpcMessageHandler) HandleMessageDeliverCommand(command command.DeliverGrpc) {
	if err := g.messageApi.DeliverMessage(command.Body, command.Protocol, command.RecipientList...); err != nil {
		logger.Errorf(nil, "[gRPC-Gateway] Fail to deliver grpc message - Protocol: [%s], RecipientList: [%s], Err: [%s]", command.Protocol, strings.Join(command.RecipientList, ", "), err.Error())
	}
}
