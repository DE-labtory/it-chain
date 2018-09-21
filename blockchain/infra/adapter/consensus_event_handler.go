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
)

type BlockApiForCommitAndStage interface {
	CommitBlock(block blockchain.DefaultBlock) error
	StageBlock(block blockchain.DefaultBlock) error
}

type ConsensusEventHandler struct {
	BlockSyncState *blockchain.BlockSyncState
	BlockApi       BlockApiForCommitAndStage
}

func NewConsensusEventHandler(blockSyncState *blockchain.BlockSyncState, blockApi BlockApiForCommitAndStage) *ConsensusEventHandler {

	return &ConsensusEventHandler{
		BlockSyncState: blockSyncState,
		BlockApi:       blockApi,
	}
}

/**
receive consensus finished event
if block sync is on progress, change state to 'staged' and add to block pool
if block sync is not on progress, commit block
*/
func (c *ConsensusEventHandler) HandleConsensusFinishedEvent(event event.ConsensusFinished) error {
	receivedBlock := extractBlockFromEvent(event)

	if len(receivedBlock.TxList) == 0 {
		return ErrBlockNil
	}

	if receivedBlock.Seal == nil {
		return ErrBlockSealNil
	}

	if c.BlockSyncState.IsProgressing() {
		c.BlockApi.StageBlock(*receivedBlock)
	} else {
		c.BlockApi.CommitBlock(*receivedBlock)
	}

	return nil
}

func extractBlockFromEvent(event event.ConsensusFinished) *blockchain.DefaultBlock {

	d := blockchain.DefaultBlock{}
	common.Deserialize(event.Body, &d)
	return &d
}
