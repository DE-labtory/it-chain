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

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

var ErrEmptyBlock = errors.New("block is nil")

type Publisher func(exchange string, topic string, data interface{}) (err error) //해당 publish함수는 midgard 에서 의존성 주입을 받기 위해 interface로 작성한다.
//모든 의존성 주입은 컴포넌트.go 에서 이루어짐

type CommandService struct {
	publisher Publisher
}

func NewCommandService(publisher Publisher) *CommandService {
	return &CommandService{
		publisher: publisher,
	}
}

func (c *CommandService) SendBlockValidateCommand(block blockchain.Block) error {
	if block == nil {
		return ErrEmptyBlock
	}

	command := blockchain.BlockValidateCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Block: block,
	}

	return c.publisher("Event", "Block", command)
}
