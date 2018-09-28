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

package api_gateway

import (
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/grpc_gateway"
	"github.com/it-chain/iLogger"
	"github.com/pkg/errors"
)

type ConnectionCommandApi struct {
}

func NewConnectionCommandApi() *ConnectionCommandApi {
	return &ConnectionCommandApi{}
}

func (cca *ConnectionCommandApi) dial(address string) (string, string, error) {
	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)
	defer client.Close()

	createConnectionCommand := command.CreateConnection{
		Address: address,
	}

	var callBackErr error
	var callBackGrpcGatewayAddress string
	var callBackConnectionId string

	iLogger.Infof(nil, "[Api_gateway] Creating connection - Address: [%s]", address)
	err := client.Call("connection.create", createConnectionCommand, func(connection grpc_gateway.Connection, err rpc.Error) {

		if !err.IsNil() {
			iLogger.Fatalf(nil, "[Api_gateway] Fail to create connection - Address: [%s]", address)
			callBackErr = errors.New(err.Message)
			return
		}

		iLogger.Infof(nil, "[Api_gateway] Connection created - gRPC-Address: [%s], Id:[%s]", connection.GrpcGatewayAddress, connection.ConnectionID)
		callBackGrpcGatewayAddress = connection.GrpcGatewayAddress
		callBackConnectionId = connection.ConnectionID
		callBackErr = nil
	})

	if err != nil {
		iLogger.Fatal(&iLogger.Fields{"err_msg": err.Error()}, "[Api_gateway] fatal err in dial cmd")
		return "", "", err
	}

	if callBackErr != nil {
		return "", "", callBackErr
	}
	return callBackGrpcGatewayAddress, callBackConnectionId, nil
}

func (cca *ConnectionCommandApi) join(address string) error {
	config := conf.GetConfiguration()
	client := rpc.NewClient(config.Engine.Amqp)
	defer client.Close()

	joinNetworkCommand := command.JoinNetwork{
		Address: address,
	}

	var callBackErr error

	iLogger.Infof(nil, "[Api_gateway] Joining network - Address: [%s]", address)
	err := client.Call("connection.join", joinNetworkCommand, func(_ struct{}, err rpc.Error) {

		if !err.IsNil() {
			iLogger.Fatalf(nil, "[Api_gateway] Fail to join network - Address: [%s], Err: [%s]", address, err.Message)
			callBackErr = errors.New(err.Message)
			return
		}

		callBackErr = nil
		iLogger.Info(nil, "[Api_gateway] Successfully request to join network")
	})

	if err != nil {
		iLogger.Fatal(&iLogger.Fields{"err_msg": err.Error()}, "[Api_gateway] fatal err in join cmd")
		return err
	}

	if callBackErr != nil {
		return callBackErr
	}

	return nil
}
