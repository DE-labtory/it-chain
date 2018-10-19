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

package mock

import "github.com/it-chain/engine/common/command"

type GrpcCall func(processId string, queue string, params interface{}, callback interface{}) error

type Client struct {
	ProcessId string
	CallFunc  GrpcCall // network manager grpc call
}

func NewClient(processId string, callFunc GrpcCall) Client {
	client := Client{
		ProcessId: processId,
		CallFunc:  callFunc,
	}
	return client
}

func (c *Client) Call(queue string, params interface{}, callback interface{}) error {
	return c.CallFunc(c.ProcessId, queue, params, callback)
}

type ConsumeFunc func(processId string, queue string, handler func(command command.ReceiveGrpc) error) error

type Server struct {
	ProcessId   string
	ConsumeFunc func(
		processId string,
		queue string,
		handler func(command command.ReceiveGrpc) error) error // network manager grpc consume
}

func NewServer(processId string, consumeFunc ConsumeFunc) Server {
	server := Server{
		ProcessId:   processId,
		ConsumeFunc: consumeFunc,
	}
	return server
}

func (s Server) Register(queue string, handler func(command command.ReceiveGrpc) error) error {
	return s.ConsumeFunc(s.ProcessId, queue, handler)
}
