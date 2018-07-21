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

package infra

import (
	"net/rpc"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"github.com/vibhavp/amqp-rpc"
)

type RpcServer struct {
	conn         amqp.Connection
	serverCodec  *rpc.ServerCodec
	server       *rpc.Server
	isServed     bool
	registerList []*interface{}
}

func NewRpcServer(conn amqp.Connection, routingKey string) (*RpcServer, error) {
	codec, err := amqprpc.NewServerCodec(&conn, routingKey, amqprpc.JSONCodec{})
	if err != nil {
		return nil, err
	}
	return &RpcServer{
		conn:         conn,
		serverCodec:  &codec,
		server:       &rpc.Server{},
		isServed:     false,
		registerList: make([]*interface{}, 0),
	}, nil
}

func (r *RpcServer) Register(instance interface{}) error {
	r.registerList = append(r.registerList, &instance)
	return r.server.Register(instance)
}

func (r *RpcServer) Serve() error {
	if len(r.registerList) == 0 {
		return errors.New("instance is not registered")
	}
	r.server.ServeCodec(r.serverCodec)
	if r.isServed {
		return errors.New("already served")
	}
	return nil
}
