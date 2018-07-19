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

	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/infra/adapter"
	"github.com/it-chain/engine/p2p/test/mock"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
)

func TestEventHandler_HandleConnCreatedEvent(t *testing.T) {

	tests := map[string]struct {
		input struct {
			nodeId  string
			address string
		}
		err error
	}{
		"success": {
			input: struct {
				nodeId  string
				address string
			}{nodeId: string("123"), address: string("123")},
			err: eventstore.ErrNilStore,
		},
		"empty address test": {
			input: struct {
				nodeId  string
				address string
			}{nodeId: string("123"), address: string("")},
			err: p2p.ErrEmptyAddress,
		},
	}

	communicationApi := &mock.MockCommunicationApi{}

	communicationApi.DeliverPLTableFunc = func(connectionId string) error {
		return nil
	}

	eventHandler := adapter.NewEventHandler(communicationApi)

	for testName, test := range tests {
		t.Logf("running test case %s", testName)
		err := eventHandler.HandleConnCreatedEvent(p2p.ConnectionCreatedEvent{EventModel: midgard.EventModel{ID: test.input.nodeId}, Address: test.input.address})
		assert.Equal(t, err, test.err)
	}

}

func TestEventHandler_HandleConnDisconnectedEvent(t *testing.T) {

	tests := map[string]struct {
		input struct {
			id string
		}
		err error
	}{
		"success": {
			input: struct {
				id string
			}{id: string(123)},
			err: nil,
		},
		"empty node id test": {
			input: struct {
				id string
			}{id: string("")},
			err: adapter.ErrEmptyPeerId,
		},
	}

	communicationApi := &mock.MockCommunicationApi{}

	eventHandler := adapter.NewEventHandler(communicationApi)

	for testName, test := range tests {

		t.Logf("running test case %s", testName)

		event := p2p.ConnectionDisconnectedEvent{
			EventModel: midgard.EventModel{
				ID: test.input.id,
			},
		}

		err := eventHandler.HandleConnDisconnectedEvent(event)

		assert.Equal(t, err, test.err)

	}
}
