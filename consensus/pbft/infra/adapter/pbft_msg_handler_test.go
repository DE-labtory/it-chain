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

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/infra/adapter"
	"github.com/it-chain/engine/consensus/pbft/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestPbftMsgHandler_HandleGrpcMsgCommand(t *testing.T) {
	proposeMsgByte, _ := common.Serialize(makeMockProposeMsg())
	prevoteMsgByte, _ := common.Serialize(makeMockPrevoteMsg())
	preCommitMsgByte, _ := common.Serialize(makeMockPreCommitMsg())

	tests := map[string]struct {
		input struct {
			cmd command.ReceiveGrpc
		}
		err error
	}{
		"ProposeMsg test": {
			input: struct {
				cmd command.ReceiveGrpc
			}{
				cmd: command.ReceiveGrpc{
					MessageId:    "MockMsg1",
					Body:         proposeMsgByte,
					ConnectionID: "connection1",
					Protocol:     "ProposeMsgProtocol",
				},
			},
			err: nil,
		},
		"PrevoteMsg test": {
			input: struct {
				cmd command.ReceiveGrpc
			}{
				cmd: command.ReceiveGrpc{
					MessageId:    "MockMsg2",
					Body:         prevoteMsgByte,
					ConnectionID: "connection1",
					Protocol:     "PrevoteMsgProtocol",
				},
			},
			err: nil,
		},
		"PreCommitMsg test": {
			input: struct {
				cmd command.ReceiveGrpc
			}{
				cmd: command.ReceiveGrpc{
					MessageId:    "MockMsg3",
					Body:         preCommitMsgByte,
					ConnectionID: "connection1",
					Protocol:     "PreCommitMsgProtocol",
				},
			},
			err: nil,
		},
		"Wrong protocol test": {
			input: struct {
				cmd command.ReceiveGrpc
			}{
				cmd: command.ReceiveGrpc{
					MessageId:    "MockMsg4",
					Body:         proposeMsgByte,
					ConnectionID: "connection1",
					Protocol:     "PrevoteMsgProtocol",
				},
			},
			err: adapter.DeserializingError,
		},
		"Undefined protocol test": {
			input: struct {
				cmd command.ReceiveGrpc
			}{
				cmd: command.ReceiveGrpc{
					MessageId:    "MockMsg5",
					Body:         proposeMsgByte,
					ConnectionID: "connection1",
					Protocol:     "UnknownProtocol",
				},
			},
			err: adapter.UndefinedProtocolError,
		},
	}

	p := adapter.NewPbftMsgHandler(newMockStateApiForPbftMsgHandler(t))

	for testName, test := range tests {
		t.Logf("running test case [%s]", testName)

		p.HandleGrpcMsgCommand(test.input.cmd)
	}
}

func newMockStateApiForPbftMsgHandler(t *testing.T) adapter.StateMsgApi {
	mockApi := &mock.StateApi{}
	mockApi.HandleProposeMsgFunc = func() error {
		if msg.SenderID == "sender1" {
			assert.NotNil(t, msg.Representative)
			assert.NotNil(t, msg.ProposedBlock)
			return nil
		}

		return errors.New("HandleProposeMsg error")
	}
	mockApi.HandlePrevoteMsgFunc = func(msg pbft.PrevoteMsg) error {
		if msg.SenderID == "sender2" {
			assert.NotNil(t, msg.BlockHash)
			return nil
		}

		return errors.New("HandlePrevoteMsg error")
	}
	mockApi.HandlePreCommitMsgFunc = func(msg pbft.PreCommitMsg) error {
		if msg.SenderID == "sender3" {
			return nil
		}

		return errors.New("HandlePreCommitMsg error")
	}

	return mockApi
}

func makeMockProposeMsg() pbft.ProposeMsg {
	return pbft.ProposeMsg{
		StateID: pbft.StateID{
			ID: "state1",
		},
		SenderID:       "sender1",
		Representative: make([]pbft.Representative, 0),
		ProposedBlock: pbft.ProposedBlock{
			Seal: []byte{'s', 'e', 'a', 'l'},
			Body: []byte{'b', 'o', 'd', 'y'},
		},
	}
}

func makeMockPrevoteMsg() pbft.PrevoteMsg {
	return pbft.PrevoteMsg{
		StateID: pbft.StateID{
			ID: "state1",
		},
		SenderID:  "sender2",
		BlockHash: []byte{'h', 'a', 's', 'h'},
	}
}

func makeMockPreCommitMsg() pbft.PreCommitMsg {
	return pbft.PreCommitMsg{
		StateID: pbft.StateID{
			ID: "state1",
		},
		SenderID: "sender3",
	}
}
