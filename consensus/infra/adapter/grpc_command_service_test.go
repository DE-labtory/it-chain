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
	"testing"

	"github.com/it-chain/engine/consensus"
	"github.com/stretchr/testify/assert"
)

func TestPropagateService_BroadcastPrePrepareMsg(t *testing.T) {
	tests := map[string]struct {
		input struct {
			msg consensus.PrePrepareMsg
		}
		err error
	}{
		"success": {
			input: struct {
				msg consensus.PrePrepareMsg
			}{
				msg: consensus.PrePrepareMsg{
					ConsensusId:    consensus.ConsensusId{"c1"},
					SenderId:       "s1",
					Representative: make([]*consensus.Representative, 0),
					ProposedBlock: consensus.ProposedBlock{
						Seal: make([]byte, 0),
						Body: make([]byte, 0),
					},
				},
			},
			err: nil,
		},
		"Consensus ID empty test": {
			input: struct {
				msg consensus.PrePrepareMsg
			}{
				msg: consensus.PrePrepareMsg{
					ConsensusId:    consensus.ConsensusId{""},
					SenderId:       "s1",
					Representative: make([]*consensus.Representative, 0),
					ProposedBlock: consensus.ProposedBlock{
						Seal: make([]byte, 0),
						Body: make([]byte, 0),
					},
				},
			},
			err: errors.New("Consensus ID is empty"),
		},
		"Block empty test": {
			input: struct {
				msg consensus.PrePrepareMsg
			}{
				msg: consensus.PrePrepareMsg{
					ConsensusId:    consensus.ConsensusId{"c1"},
					SenderId:       "s1",
					Representative: make([]*consensus.Representative, 0),
					ProposedBlock: consensus.ProposedBlock{
						Seal: make([]byte, 0),
						Body: nil,
					},
				},
			},
			err: errors.New("Block is empty"),
		},
	}

	publish := func(exchange string, topic string, data interface{}) (e error) {
		assert.Equal(t, "Command", exchange)
		assert.Equal(t, "message.deliver", topic)

		return nil
	}

	representatives := make([]*consensus.Representative, 0)
	propagateService := NewGrpcCommandService(publish, representatives)

	for testName, test := range tests {
		t.Logf("running test case [%s]", testName)

		err := propagateService.BroadcastPrePrepareMsg(test.input.msg)

		assert.Equal(t, test.err, err)
	}
}

func TestPropagateService_BroadcastPrepareMsg(t *testing.T) {
	tests := map[string]struct {
		input struct {
			msg consensus.PrepareMsg
		}
		err error
	}{
		"success": {
			input: struct {
				msg consensus.PrepareMsg
			}{
				msg: consensus.PrepareMsg{
					ConsensusId: consensus.ConsensusId{"c1"},
					SenderId:    "s1",
					BlockHash:   make([]byte, 0),
				},
			},
			err: nil,
		},
		"Consensus ID empty test": {
			input: struct {
				msg consensus.PrepareMsg
			}{
				msg: consensus.PrepareMsg{
					ConsensusId: consensus.ConsensusId{""},
					SenderId:    "s1",
					BlockHash:   make([]byte, 0),
				},
			},
			err: errors.New("Consensus ID is empty"),
		},
		"Block hash empty test": {
			input: struct {
				msg consensus.PrepareMsg
			}{
				msg: consensus.PrepareMsg{
					ConsensusId: consensus.ConsensusId{"c1"},
					SenderId:    "s1",
					BlockHash:   nil,
				},
			},
			err: errors.New("Block hash is empty"),
		},
	}

	publish := func(exchange string, topic string, data interface{}) (e error) {
		assert.Equal(t, "Command", exchange)
		assert.Equal(t, "message.deliver", topic)

		return nil
	}

	representatives := make([]*consensus.Representative, 0)
	propagateService := NewGrpcCommandService(publish, representatives)

	for testName, test := range tests {
		t.Logf("running test case [%s]", testName)

		err := propagateService.BroadcastPrepareMsg(test.input.msg)

		assert.Equal(t, test.err, err)
	}
}

func TestPropagateService_BroadcastCommitMsg(t *testing.T) {
	tests := map[string]struct {
		input struct {
			msg consensus.CommitMsg
		}
		err error
	}{
		"success": {
			input: struct {
				msg consensus.CommitMsg
			}{
				msg: consensus.CommitMsg{
					ConsensusId: consensus.ConsensusId{"c1"},
					SenderId:    "s1",
				},
			},
			err: nil,
		},
		"Consensus ID empty test": {
			input: struct {
				msg consensus.CommitMsg
			}{
				msg: consensus.CommitMsg{
					ConsensusId: consensus.ConsensusId{""},
					SenderId:    "s1",
				},
			},
			err: errors.New("Consensus ID is empty"),
		},
	}

	publish := func(exchange string, topic string, data interface{}) (e error) {
		assert.Equal(t, "Command", exchange)
		assert.Equal(t, "message.deliver", topic)

		return nil
	}

	representatives := make([]*consensus.Representative, 0)
	propagateService := NewGrpcCommandService(publish, representatives)

	for testName, test := range tests {
		t.Logf("running test case [%s]", testName)

		err := propagateService.BroadcastCommitMsg(test.input.msg)

		assert.Equal(t, test.err, err)
	}
}
