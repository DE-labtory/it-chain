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
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/consensus/pbft"
)

type Publish func(topic string, data interface{}) (err error)

type EventService struct {
	publish Publish
}

func NewEventService(publish Publish) EventService {
	return EventService{
		publish: publish,
	}
}

func (es EventService) ConfirmBlock(block pbft.ProposedBlock) error {

	// todo : block을 consensus finished event로 변경하여 날려야함
	e := event.ConsensusFinished{}

	return es.publish("block.confirm", e)
}
