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
	"time"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/blockchain/infra/mem"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/it-chain/engine/common"
	event2 "github.com/it-chain/engine/common/event"
	"github.com/magiconair/properties/assert"
)

func TestNewConsensusEventHandler(t *testing.T) {
	syncStateRepository := mem.NewSyncStateRepository()
	blockApi := &mock.BlockApi{}

	consensusEventHandler := &adapter.ConsensusEventHandler{SyncStateRepository: syncStateRepository, BlockApi: blockApi}
	generated := adapter.NewConsensusEventHandler(syncStateRepository, blockApi)

	assert.Equal(t, consensusEventHandler, generated)
}

func TestConsensusEventHandler_HandleConsensusFinishedEvent(t *testing.T) {

	block1 := blockchain.DefaultBlock{
		Seal:      []byte{'s', 'e', 'a', 'l'},
		PrevSeal:  []byte{},
		Height:    0,
		TxList:    []*blockchain.DefaultTransaction{},
		Timestamp: time.Time{},
		Creator:   "",
		State:     "",
	}

	body1, _ := common.Serialize(block1)

	block2 := blockchain.DefaultBlock{
		Seal:      nil,
		PrevSeal:  []byte{},
		Height:    0,
		TxList:    []*blockchain.DefaultTransaction{},
		Timestamp: time.Time{},
		Creator:   "",
		State:     "",
	}

	body2, _ := common.Serialize(block2)

	event1 := event2.ConsensusFinished{
		Seal: []byte{'s', 'e', 'a', 'l'},
		Body: body1,
	}

	event2 := event2.ConsensusFinished{
		Seal: nil,
		Body: body2,
	}

	syncStateRepository := mem.NewSyncStateRepository()
	consensusEventHandler := SetConsensusEventHandler(syncStateRepository)

	//	when
	err1 := consensusEventHandler.HandleConsensusFinishedEvent(event1)

	//then
	assert.Equal(t, err1, nil)

	//	when
	err2 := consensusEventHandler.HandleConsensusFinishedEvent(event2)

	//	then
	assert.Equal(t, err2, adapter.ErrBlockSealNil)

}

func SetConsensusEventHandler(state blockchain.SyncStateRepository) *adapter.ConsensusEventHandler {
	blockApi := &mock.BlockApi{}
	blockApi.CommitBlockFunc = func(block blockchain.DefaultBlock) error {
		return nil
	}
	blockApi.StageBlockFunc = func(block blockchain.DefaultBlock) {
		return
	}

	return &adapter.ConsensusEventHandler{
		SyncStateRepository: state,
		BlockApi:            blockApi,
	}
}
