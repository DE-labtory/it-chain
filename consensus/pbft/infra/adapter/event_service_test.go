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
	"testing"

	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/infra/adapter"
	"github.com/it-chain/engine/consensus/pbft/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestEventService_ConfirmBlock(t *testing.T) {
	mockEventService := mock.EventService{}
	mockEventService.PublishFunc = func(topic string, event interface{}) error {
		// todo : 생성된 event 비교해야함
		assert.Equal(t, "block.confirm", topic)
		return nil
	}

	eventService := adapter.NewEventService(mockEventService.PublishFunc)

	block := pbft.ProposedBlock{
		Seal: []byte("seal"),
		Body: []byte("body"),
	}

	err := eventService.ConfirmBlock(block)
	assert.Nil(t, err)
}
