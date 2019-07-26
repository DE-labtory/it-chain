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

package api

import (
	"strings"

	"github.com/DE-labtory/iLogger"
	"github.com/DE-labtory/it-chain/grpc_gateway"
)

type MessageApi struct {
	grpcService grpc_gateway.GrpcService
}

func NewMessageApi(grpcService grpc_gateway.GrpcService) *MessageApi {
	return &MessageApi{
		grpcService: grpcService,
	}
}

func (c MessageApi) DeliverMessage(body []byte, protocol string, ids ...string) error {
	iLogger.Infof(nil, "[gRPC-Gateway] Sending gRPC message... - Protocol: [%s], RecipientList: [%s]", protocol, strings.Join(ids, ", "))
	return c.grpcService.SendMessages(body, protocol, ids...)
}
