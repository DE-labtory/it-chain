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

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/blockchain/test/mock"
	event2 "github.com/it-chain/engine/common/event"
	"github.com/magiconair/properties/assert"
)

func TestNewConsensusEventHandler(t *testing.T) {
	blockSyncState := blockchain.NewBlockSyncState()
	blockApi := &mock.BlockApi{}

	consensusEventHandler := &adapter.ConsensusEventHandler{BlockSyncState: blockSyncState, BlockApi: blockApi}
	generated := adapter.NewConsensusEventHandler(blockSyncState, blockApi)

	assert.Equal(t, consensusEventHandler, generated)
}

func TestConsensusEventHandler_HandleConsensusFinishedEvent(t *testing.T) {

	event1 := event2.ConsensusFinished{
		PrevSeal: []byte{},
		Height:   12,
		TxList: []event2.Tx{event2.Tx{
			ID: "1",
		}},
		Creator: []byte{},
	}

	event2 := event2.ConsensusFinished{
		PrevSeal: []byte{},
		Height:   12,
		TxList:   []event2.Tx{},
		Creator:  []byte{},
	}

	consensusEventHandler := SetConsensusEventHandler(&blockchain.BlockSyncState{Id: "1", IsProgress: true})

	//	when
	consensusEventHandler.HandleConsensusFinishedEvent(event1)

	//then
	assert.Equal(t, consensusEventHandler.HandleConsensusFinishedEvent(event1), nil)

	//	when
	consensusEventHandler.HandleConsensusFinishedEvent(event2)

	//	then
	assert.Equal(t, consensusEventHandler.HandleConsensusFinishedEvent(event2), adapter.ErrBlockNil)

}

func SetConsensusEventHandler(state *blockchain.BlockSyncState) *adapter.ConsensusEventHandler {
	blockApi := &mock.BlockApi{}
	blockApi.CommitBlockFunc = func(block blockchain.DefaultBlock) error {
		return nil
	}
	blockApi.StageBlockFunc = func(block blockchain.DefaultBlock) error {
		return nil
	}

	return &adapter.ConsensusEventHandler{
		BlockSyncState: state,
		BlockApi:       blockApi,
	}
}
