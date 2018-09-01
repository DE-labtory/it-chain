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
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/blockchain/api"
)

type ConsensusEventHandler struct {
	blockSyncState *blockchain.BlockSyncState
	blockPool      *blockchain.BlockPool
	blockApi       *api.BlockApi
}

func NewConsensusEventHandler(blockApi *api.BlockApi) *ConsensusEventHandler {

	return &ConsensusEventHandler{}
}

/**
receive consensus finished event
if block sync is on progress, change state to 'staged' and add to block pool
if block sync is not on progress, commit block
*/
func (c *ConsensusEventHandler) HandlerConsensusFinishedEvent(event event.ConsensusFinished) error {

	receivedBlock := extractBlockFromEvent(event)

	if c.blockSyncState.IsProgressing() {
		receivedBlock.SetState(blockchain.Staged)
		c.blockPool.Add(receivedBlock)

		return nil
	} else {
		c.blockApi.CommitBlock(*receivedBlock)
	}

	return nil
}

func extractBlockFromEvent(event event.ConsensusFinished) *blockchain.DefaultBlock{
	txList := blockchain.ConvertToTransactionList(event.TxList)

	return &blockchain.DefaultBlock{
		PrevSeal:event.PrevSeal,
		Height:event.Height,
		TxList:txList,
		Creator:event.Creator,
	}
}
