/*
 * Copyright 2018 DE-labtory
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

import "github.com/DE-labtory/engine/grpc_gateway"

type GrpcService struct {
	DialFunc                func(address string) (grpc_gateway.Connection, error)
	CloseConnectionFunc     func(connID string) error
	SendMessagesFunc        func(message []byte, protocol string, connIDs ...string) error
	GetAllConnectionsFunc   func() ([]grpc_gateway.Connection, error)
	CloseAllConnectionsFunc func() error
	IsConnectionExistFunc   func(connectionID string) bool
	GetHostIDFunc           func() string
}

func (g GrpcService) Dial(address string) (grpc_gateway.Connection, error) {
	return g.DialFunc(address)
}

func (g GrpcService) CloseConnection(connID string) error {
	return g.CloseConnectionFunc(connID)
}

func (g GrpcService) SendMessages(message []byte, protocol string, connIDs ...string) error {
	return g.SendMessagesFunc(message, protocol, connIDs...)
}

func (g GrpcService) GetAllConnections() ([]grpc_gateway.Connection, error) {
	return g.GetAllConnectionsFunc()
}

func (g GrpcService) CloseAllConnections() error {
	return g.CloseAllConnectionsFunc()
}

func (g GrpcService) IsConnectionExist(connectionID string) bool {
	return g.IsConnectionExistFunc(connectionID)
}

func (g GrpcService) GetHostID() string {
	return g.GetHostIDFunc()
}
