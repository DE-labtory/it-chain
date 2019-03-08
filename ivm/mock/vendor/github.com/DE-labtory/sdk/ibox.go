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
 *
 */

package sdk

import (
	"os"

	"github.com/DE-labtory/sdk/logger"
	"github.com/DE-labtory/sdk/pb"
)

type IBox struct {
	port    int
	handler RequestHandler
}

func NewIBox(port int) *IBox {
	return &IBox{
		port:    port,
		handler: nil,
	}
}

func (i *IBox) SetHandler(handler RequestHandler) {
	i.handler = handler
}

func (i *IBox) On(timeout int) error {
	if i.handler == nil {
		logger.Panic(nil, "handler is nil")
		os.Exit(1)
	}
	server := NewServer(i.port)
	cell := NewCell(i.handler.Name())

	serverHandler := func(request *pb.Request) *pb.Response {
		return i.handler.Handle(request, cell)
	}
	server.SetHandler(serverHandler)
	return server.Listen(timeout)
}
