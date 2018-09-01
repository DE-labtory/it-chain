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
	"testing"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/test/mock"
	event2 "github.com/it-chain/engine/common/event"
	"github.com/magiconair/properties/assert"
)

func TestNewConsensusEventHandler(t *testing.T) {
	blockSyncState := blockchain.NewBlockSyncState()
	blockPool := blockchain.NewBlockPool()
	blockApi := &mock.BlockApi{}

	consensusEventHandler := &ConsensusEventHandler{blockSyncState: blockSyncState, blockPool: blockPool, blockApi: blockApi}
	generated := NewConsensusEventHandler(blockSyncState, blockPool, blockApi)

	assert.Equal(t, consensusEventHandler, generated)
}

func TestConsensusEventHandler_HandleConsensusFinishedEvent(t *testing.T) {

	event := &event2.ConsensusFinished{
		PrevSeal: []byte{},
		Height:   12,
		TxList: []event2.Tx{event2.Tx{
			ID: "1",
		}},
		Creator: []byte{},
	}

	//	when
	consensusEventHandler1 := SetConsensusEventHandler(&blockchain.BlockSyncState{Id: "1", IsProgress: true})
	consensusEventHandler1.HandleConsensusFinishedEvent(*event)

	//then
	assert.Equal(t, consensusEventHandler1.blockPool.Get(12).GetHeight(), uint64(12))

	//	when
	consensusEventHandler2 := SetConsensusEventHandler(&blockchain.BlockSyncState{Id: "1", IsProgress: false})
	consensusEventHandler2.HandleConsensusFinishedEvent(*event)

	//	then
	assert.Equal(t, consensusEventHandler2.blockPool.Get(12), nil)

}

func SetConsensusEventHandler(state *blockchain.BlockSyncState) *ConsensusEventHandler {
	blockPool := blockchain.NewBlockPool()
	blockApi := &mock.BlockApi{}
	blockApi.CommitBlockFunc = func(block blockchain.DefaultBlock) error {
		return nil
	}

	return &ConsensusEventHandler{
		blockSyncState: state,
		blockPool:      blockPool,
		blockApi:       blockApi,
	}
}
