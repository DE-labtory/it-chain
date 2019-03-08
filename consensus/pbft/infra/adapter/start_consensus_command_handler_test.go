/*
 * Copyright 2018 DE-labtory
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

	"github.com/DE-labtory/engine/consensus/pbft"
	"github.com/DE-labtory/engine/consensus/pbft/infra/adapter"
)

func TestStartConsensusCommandHandler_HandleStartConsensusCommand(t *testing.T) {
	//mockStateApi := newMockStateApi(nil)
	//testHandler := adapter.NewStartConsensusCommandHandler(mockStateApi)
	//
	//// case 1 : success
	//expectedSeal := []byte{'s', 'e', 'a', 'l'}
	//expectedTxList := []event.Tx{}
	//for i := 0; i < 5; i++ {
	//	expectedTxList = append(expectedTxList, event.Tx{
	//		ID: string(i),
	//	})
	//}
	//
	//expectedCommand := event.BlockCreated{
	//	Seal:   expectedSeal,
	//	TxList: expectedTxList,
	//}
	//
	//testResult, testErr := testHandler.HandleStartConsensusCommand(expectedCommand)
	//
	//assert.True(t, testResult)
	//assert.Equal(t, "", testErr.Message)
	//
	//// case 2 : consensus on error
	//consensusStartError := errors.New("on consensus failed!")
	//mockStateApi = newMockStateApi(consensusStartError)
	//testHandler = adapter.NewStartConsensusCommandHandler(mockStateApi)
	//
	//testResult, testErr = testHandler.HandleStartConsensusCommand(expectedCommand)
	//
	//assert.False(t, testResult)
	//assert.Equal(t, consensusStartError.Error(), testErr.Message)
}

func newMockStateApi(err error) adapter.StateStartApi {
	mockStateApi := adapter.StateStartApi{}

	mockStateApi.StartConsensus = func(proposedBlock pbft.ProposedBlock) error {
		return err
	}

	return mockStateApi
}
