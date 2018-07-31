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
	"errors"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

//kind of error
var ErrEmptyPeerId = errors.New("empty nodeid proposed")
var ErrEmptyLeaderId = errors.New("empty leader id proposed")
var ErrEmptyConnectionId = errors.New("empty connection ")

type Publish func(exchange string, topic string, data interface{}) (err error) // 나중에 의존성 주입을 해준다.

// message dispatcher sends messages to other nodes in p2p network
type GrpcCommandService struct {
	publish Publish // midgard.client.Publish
}

func NewGrpcCommandService(publish Publish) *GrpcCommandService {

	return &GrpcCommandService{

		publish: publish,
	}
}

func CreateGrpcDeliverCommand(protocol string, body interface{}) (command.DeliverGrpc, error) {

	data, err := common.Serialize(body)

	if err != nil {

		return command.DeliverGrpc{}, err
	}

	return command.DeliverGrpc{

		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		RecipientList: make([]string, 0),
		Body:          data,
		Protocol:      protocol,
	}, err
}
