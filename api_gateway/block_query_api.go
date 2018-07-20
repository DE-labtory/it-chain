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

package api_gateway

import (
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/yggdrasill"

	"sync"
)

type BlockQueryApi struct {
	blockPoolRepository     BlockPoolRepository
	commitedBlockRepository CommitedBlockRepository
}

type BlockPoolRepository interface {
	AddCreatedBlock(block blockchain.DefaultBlock)
	GetStagedBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
	GetStagedBlockById(id string) (blockchain.DefaultBlock, error)
	GetFirstStagedBlock() (blockchain.DefaultBlock, error)
}

type BlockPoolRepositoryImpl struct {
	Blocks []blockchain.Block
}

type CommitedBlockRepository interface {
	Save(block blockchain.DefaultBlock) error
	GetLastBlock() (blockchain.DefaultBlock, error)
	GetBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
}

type CommitedBlockRepositoryImpl struct {
	mux *sync.RWMutex
	yggdrasill.BlockStorageManager
}
