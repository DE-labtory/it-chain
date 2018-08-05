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

package mock

import (
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common/event"
)

type BlockApi struct {
	AddBlockToPoolFunc            func(block blockchain.Block) error
	CheckAndSaveBlockFromPoolFunc func(height blockchain.BlockHeight) error
	CommitProposedBlockFunc       func(txList []*blockchain.DefaultTransaction) error
}

func (api BlockApi) AddBlockToPool(block blockchain.Block) error {
	return api.AddBlockToPoolFunc(block)
}

func (api BlockApi) CheckAndSaveBlockFromPool(height blockchain.BlockHeight) error {
	return api.CheckAndSaveBlockFromPoolFunc(height)
}

func (api BlockApi) CommitProposedBlock(txList []*blockchain.DefaultTransaction) error {
	return api.CommitProposedBlockFunc(txList)
}

type MockSyncBlockApi struct {
	SyncedCheckFunc func(block blockchain.Block) error
}

func (ba MockSyncBlockApi) SyncedCheck(block blockchain.Block) error {
	return ba.SyncedCheckFunc(block)
}

type CommitEventHandler struct {
	HandleFunc func(event event.BlockCommitted)
}