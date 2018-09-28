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

package pbft_test

import (
	"errors"
	"testing"

	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestPropagateService_BroadcastProposeMsg(t *testing.T) {
	tests := map[string]struct {
		input struct {
			msg pbft.ProposeMsg
		}
		err error
	}{
		"success": {
			input: struct {
				msg pbft.ProposeMsg
			}{
				msg: pbft.ProposeMsg{
					StateID:        pbft.StateID{"c1"},
					SenderID:       "s1",
					Representative: make([]*pbft.Representative, 0),
					ProposedBlock: pbft.ProposedBlock{
						Seal: make([]byte, 0),
						Body: make([]byte, 0),
					},
				},
			},
			err: nil,
		},
		"State ID empty test": {
			input: struct {
				msg pbft.ProposeMsg
			}{
				msg: pbft.ProposeMsg{
					StateID:        pbft.StateID{""},
					SenderID:       "s1",
					Representative: make([]*pbft.Representative, 0),
					ProposedBlock: pbft.ProposedBlock{
						Seal: make([]byte, 0),
						Body: make([]byte, 0),
					},
				},
			},
			err: errors.New("State ID is empty"),
		},
		"Block empty test": {
			input: struct {
				msg pbft.ProposeMsg
			}{
				msg: pbft.ProposeMsg{
					StateID:        pbft.StateID{"c1"},
					SenderID:       "s1",
					Representative: make([]*pbft.Representative, 0),
					ProposedBlock: pbft.ProposedBlock{
						Seal: make([]byte, 0),
						Body: nil,
					},
				},
			},
			err: errors.New("Block is empty"),
		},
	}

	mockEventService := mock.EventService{}
	mockEventService.PublishFunc = func(topic string, event interface{}) error {
		assert.Equal(t, "message.deliver", topic)

		return nil
	}

	representatives := make([]*pbft.Representative, 0)
	propagateService := pbft.NewPropagateService(mockEventService)

	for testName, test := range tests {
		t.Logf("running test case [%s]", testName)

		err := propagateService.BroadcastProposeMsg(test.input.msg, representatives)

		assert.Equal(t, test.err, err)
	}
}

func TestPropagateService_BroadcastPrevoteMsg(t *testing.T) {
	tests := map[string]struct {
		input struct {
			msg pbft.PrevoteMsg
		}
		err error
	}{
		"success": {
			input: struct {
				msg pbft.PrevoteMsg
			}{
				msg: pbft.PrevoteMsg{
					StateID:   pbft.StateID{"c1"},
					SenderID:  "s1",
					BlockHash: make([]byte, 0),
				},
			},
			err: nil,
		},
		"State ID empty test": {
			input: struct {
				msg pbft.PrevoteMsg
			}{
				msg: pbft.PrevoteMsg{
					StateID:   pbft.StateID{""},
					SenderID:  "s1",
					BlockHash: make([]byte, 0),
				},
			},
			err: errors.New("State ID is empty"),
		},
		"Block hash empty test": {
			input: struct {
				msg pbft.PrevoteMsg
			}{
				msg: pbft.PrevoteMsg{
					StateID:   pbft.StateID{"c1"},
					SenderID:  "s1",
					BlockHash: nil,
				},
			},
			err: errors.New("Block hash is empty"),
		},
	}

	mockEventService := mock.EventService{}
	mockEventService.PublishFunc = func(topic string, event interface{}) error {
		assert.Equal(t, "message.deliver", topic)

		return nil
	}

	representatives := make([]*pbft.Representative, 0)
	propagateService := pbft.NewPropagateService(mockEventService)

	for testName, test := range tests {
		t.Logf("running test case [%s]", testName)

		err := propagateService.BroadcastPrevoteMsg(test.input.msg, representatives)

		assert.Equal(t, test.err, err)
	}
}

func TestPropagateService_BroadcastPreCommitMsg(t *testing.T) {
	tests := map[string]struct {
		input struct {
			msg pbft.PreCommitMsg
		}
		err error
	}{
		"success": {
			input: struct {
				msg pbft.PreCommitMsg
			}{
				msg: pbft.PreCommitMsg{
					StateID:  pbft.StateID{"c1"},
					SenderID: "s1",
				},
			},
			err: nil,
		},
		"State ID empty test": {
			input: struct {
				msg pbft.PreCommitMsg
			}{
				msg: pbft.PreCommitMsg{
					StateID:  pbft.StateID{""},
					SenderID: "s1",
				},
			},
			err: errors.New("State ID is empty"),
		},
	}

	mockEventService := mock.EventService{}
	mockEventService.PublishFunc = func(topic string, event interface{}) error {
		assert.Equal(t, "message.deliver", topic)

		return nil
	}

	representatives := make([]*pbft.Representative, 0)
	propagateService := pbft.NewPropagateService(mockEventService)

	for testName, test := range tests {
		t.Logf("running test case [%s]", testName)

		err := propagateService.BroadcastPreCommitMsg(test.input.msg, representatives)

		assert.Equal(t, test.err, err)
	}
}
