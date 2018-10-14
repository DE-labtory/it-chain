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
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/sdk/logger"
)

type BlockApiForCommitAndStage interface {
	CommitBlock(block blockchain.DefaultBlock) error
	StageBlock(block blockchain.DefaultBlock)
}

type ConsensusEventHandler struct {
	SyncStateRepository blockchain.SyncStateRepository
	BlockApi            BlockApiForCommitAndStage
}

func NewConsensusEventHandler(syncStateRepository blockchain.SyncStateRepository, blockApi BlockApiForCommitAndStage) *ConsensusEventHandler {

	return &ConsensusEventHandler{
		SyncStateRepository: syncStateRepository,
		BlockApi:            blockApi,
	}
}

/**
receive consensus finished event
if block sync is on progress, change state to 'staged' and add to block pool
if block sync is not on progress, commit block
*/
func (c *ConsensusEventHandler) HandleConsensusFinishedEvent(event event.ConsensusFinished) error {
	receivedBlock := extractBlockFromEvent(event)

	if receivedBlock.Seal == nil {
		return ErrBlockSealNil
	}

	syncState := c.SyncStateRepository.Get()

	if syncState.SyncProgressing {
		c.BlockApi.StageBlock(*receivedBlock)
	} else {
		if err := c.BlockApi.CommitBlock(*receivedBlock); err != nil {
			return err
		}
	}

	return nil
}

func extractBlockFromEvent(event event.ConsensusFinished) *blockchain.DefaultBlock {
	block := &blockchain.DefaultBlock{}

	if err := common.Deserialize(event.Body, block); err != nil {
		logger.Error(nil, "[Blockchain] Deserialize error")
	}

	return block
}
