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
	SyncedCheckFunc               func(block blockchain.Block) error
	AddBlockToPoolFunc            func(block blockchain.Block) error
	CheckAndSaveBlockFromPoolFunc func(height blockchain.BlockHeight) error
	SyncIsProgressingFunc         func() blockchain.ProgressState
	CommitGenesisBlockFunc        func(GenesisConfPath string) error
	CommitBlockFunc               func(block blockchain.DefaultBlock) error
	StageBlockFunc                func(block blockchain.DefaultBlock) error
	CreateProposedBlockFunc       func(txList []*blockchain.DefaultTransaction) (blockchain.DefaultBlock, error)
}

func (api BlockApi) SyncedCheck(block blockchain.Block) error {
	return api.SyncedCheckFunc(block)
}

func (api BlockApi) AddBlockToPool(block blockchain.Block) error {
	return api.AddBlockToPoolFunc(block)
}

func (api BlockApi) CheckAndSaveBlockFromPool(height blockchain.BlockHeight) error {
	return api.CheckAndSaveBlockFromPoolFunc(height)
}

func (api BlockApi) SyncIsProgressing() blockchain.ProgressState {
	return api.SyncIsProgressingFunc()
}

func (api BlockApi) CommitGenesisBlock(GenesisConfPath string) error {
	return api.CommitGenesisBlockFunc(GenesisConfPath)
}

func (api BlockApi) CommitBlock(block blockchain.DefaultBlock) error {
	return api.CommitBlockFunc(block)
}

func (api BlockApi) StageBlock(block blockchain.DefaultBlock) error {
	return api.StageBlockFunc(block)
}

func (api BlockApi) CreateProposedBlock(txList []*blockchain.DefaultTransaction) (blockchain.DefaultBlock, error) {
	return api.CreateProposedBlockFunc(txList)
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

func (h *CommitEventHandler) Handle(event event.BlockCommitted) {
	h.HandleFunc(event)
}
