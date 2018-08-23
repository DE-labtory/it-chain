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

	"errors"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/api"
	"github.com/it-chain/engine/consensus/pbft/infra/adapter"
	"github.com/it-chain/engine/consensus/pbft/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestStartConsensusCommandHandler_HandleStartConsensusCommand(t *testing.T) {
	mockStateApi := newMockStateApi(nil)
	testHandler := adapter.NewStartConsensusCommandHandler(mockStateApi)

	// case 1 : success
	expectedSeal := []byte{'s', 'e', 'a', 'l'}
	expectedTxList := []command.Tx{}
	for i := 0; i < 5; i++ {
		expectedTxList = append(expectedTxList, command.Tx{
			ID: string(i),
		})
	}

	expectedCommand := command.StartConsensus{
		Seal:   expectedSeal,
		TxList: expectedTxList,
	}

	testResult, testErr := testHandler.HandleStartConsensusCommand(expectedCommand)

	assert.True(t, testResult)
	assert.Equal(t, "", testErr.Message)

	// case 2 : consensus start error
	consensusStartError := errors.New("start consensus failed!")
	mockStateApi = newMockStateApi(consensusStartError)
	testHandler = adapter.NewStartConsensusCommandHandler(mockStateApi)

	testResult, testErr = testHandler.HandleStartConsensusCommand(expectedCommand)

	assert.False(t, testResult)
	assert.Equal(t, consensusStartError.Error(), testErr.Message)
}

func newMockStateApi(err error) api.StateApi {
	mockStateApi := mock.MockStateApi{}

	mockStateApi.StartConsensusFunc = func(proposedBlock pbft.ProposedBlock) error {
		return err
	}
	mockStateApi.HandlePrePrepareMsgFunc = func(msg pbft.PrePrepareMsg) error {
		return nil
	}
	mockStateApi.HandlePrepareMsgFunc = func(msg pbft.PrepareMsg) error {
		return nil
	}
	mockStateApi.HandleCommitMsgFunc = func(msg pbft.CommitMsg) error {
		return nil
	}

	return &mockStateApi
}
