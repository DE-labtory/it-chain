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
	"encoding/json"
	"sync"
	"testing"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/grpc_gateway"
	"github.com/it-chain/engine/grpc_gateway/infra/adapter"
	"github.com/it-chain/engine/grpc_gateway/mock"
	"github.com/magiconair/properties/assert"
)

func TestGrpcMessageHandler_HandleMessageDeliverCommand(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	deliverApi := &mock.MessageDeliverApi{}
	deliverApi.DeliverMessageFunc = func(body []byte, protocol string, ids ...string) error {
		assert.Equal(t, body, []byte("hello"))
		assert.Equal(t, protocol, "RequestPeerList")
		assert.Equal(t, ids, []string{"peer1", "peer2"})
		wg.Done()
		return nil
	}

	handlerApi := &mock.GrpcMessageHandlerApi{}
	g := adapter.NewGrpcMessageHandler(handlerApi, deliverApi)

	deliverCommand := command.DeliverGrpc{
		Body:          []byte("hello"),
		RecipientList: []string{"peer1", "peer2"},
		Protocol:      "RequestPeerList",
	}

	g.HandleMessageDeliverCommand(deliverCommand)
	wg.Wait()
}

func TestGrpcMessageHandler_HandleMessageReceiveCommand(t *testing.T) {

	connectionList := []grpc_gateway.Connection{grpc_gateway.Connection{ConnectionID: "123", GrpcGatewayAddress: "127.0.0.1:3333", ApiGatewayAddress: "127.0.0.1:3334"}}

	//given
	deliverApi := &mock.MessageDeliverApi{}
	handlerApi := &mock.GrpcMessageHandlerApi{}
	handlerApi.DialConnectionListFunc = func(connectionList []grpc_gateway.Connection) {
		//then
		assert.Equal(t, connectionList, connectionList)
	}
	handlerApi.HandleRequestPeerListFunc = func(connectionID string) {

		//then
		assert.Equal(t, connectionID, "peer1")
	}
	g := adapter.NewGrpcMessageHandler(handlerApi, deliverApi)

	peerRequestProtocolCommand := command.ReceiveGrpc{
		Protocol:     grpc_gateway.RequestPeerList,
		ConnectionID: "peer1",
	}

	b, _ := json.Marshal(connectionList)
	peerResponseProtocolCommand := command.ReceiveGrpc{
		Protocol: grpc_gateway.ResponsePeerList,
		Body:     b,
	}

	//when

	// 1. RequestPeerList protocol
	g.HandleMessageReceiveCommand(peerRequestProtocolCommand)

	// 2. ResponsePeerList protocol
	g.HandleMessageReceiveCommand(peerResponseProtocolCommand)
}
