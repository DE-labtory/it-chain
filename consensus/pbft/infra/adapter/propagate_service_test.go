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
	"errors"
	"testing"

	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/infra/adapter"
	"github.com/stretchr/testify/assert"
)

func TestPropagateService_BroadcastPrePrepareMsg(t *testing.T) {
	tests := map[string]struct {
		input struct {
			msg pbft.PrePrepareMsg
		}
		err error
	}{
		"success": {
			input: struct {
				msg pbft.PrePrepareMsg
			}{
				msg: pbft.PrePrepareMsg{
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
				msg pbft.PrePrepareMsg
			}{
				msg: pbft.PrePrepareMsg{
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
				msg pbft.PrePrepareMsg
			}{
				msg: pbft.PrePrepareMsg{
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

	representatives := make([]*pbft.Representative, 0)
	propagateService := adapter.NewPropagateService()

	for testName, test := range tests {
		t.Logf("running test case [%s]", testName)

		err := propagateService.BroadcastPrePrepareMsg(test.input.msg, representatives)

		assert.Equal(t, test.err, err)
	}
}

func TestPropagateService_BroadcastPrepareMsg(t *testing.T) {
	tests := map[string]struct {
		input struct {
			msg pbft.PrepareMsg
		}
		err error
	}{
		"success": {
			input: struct {
				msg pbft.PrepareMsg
			}{
				msg: pbft.PrepareMsg{
					StateID:   pbft.StateID{"c1"},
					SenderID:  "s1",
					BlockHash: make([]byte, 0),
				},
			},
			err: nil,
		},
		"State ID empty test": {
			input: struct {
				msg pbft.PrepareMsg
			}{
				msg: pbft.PrepareMsg{
					StateID:   pbft.StateID{""},
					SenderID:  "s1",
					BlockHash: make([]byte, 0),
				},
			},
			err: errors.New("State ID is empty"),
		},
		"Block hash empty test": {
			input: struct {
				msg pbft.PrepareMsg
			}{
				msg: pbft.PrepareMsg{
					StateID:   pbft.StateID{"c1"},
					SenderID:  "s1",
					BlockHash: nil,
				},
			},
			err: errors.New("Block hash is empty"),
		},
	}

	representatives := make([]*pbft.Representative, 0)
	propagateService := adapter.NewPropagateService()

	for testName, test := range tests {
		t.Logf("running test case [%s]", testName)

		err := propagateService.BroadcastPrepareMsg(test.input.msg, representatives)

		assert.Equal(t, test.err, err)
	}
}

func TestPropagateService_BroadcastCommitMsg(t *testing.T) {
	tests := map[string]struct {
		input struct {
			msg pbft.CommitMsg
		}
		err error
	}{
		"success": {
			input: struct {
				msg pbft.CommitMsg
			}{
				msg: pbft.CommitMsg{
					StateID:  pbft.StateID{"c1"},
					SenderID: "s1",
				},
			},
			err: nil,
		},
		"State ID empty test": {
			input: struct {
				msg pbft.CommitMsg
			}{
				msg: pbft.CommitMsg{
					StateID:  pbft.StateID{""},
					SenderID: "s1",
				},
			},
			err: errors.New("State ID is empty"),
		},
	}

	representatives := make([]*pbft.Representative, 0)
	propagateService := adapter.NewPropagateService()

	for testName, test := range tests {
		t.Logf("running test case [%s]", testName)

		err := propagateService.BroadcastCommitMsg(test.input.msg, representatives)

		assert.Equal(t, test.err, err)
	}
}
