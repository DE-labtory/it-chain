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

package mem

import (
	"sync"

	"github.com/it-chain/engine/blockchain"
	"github.com/gogo/protobuf/sortkeys"
)

type BlockPool struct {
	blockMap blockchain.BlockMap
	mux      sync.Mutex
}

func NewBlockPool() *BlockPool {
	return &BlockPool{
		blockMap: make(map[uint64]blockchain.DefaultBlock),
		mux:      sync.Mutex{},
	}
}

func (b *BlockPool) Add(block blockchain.DefaultBlock) {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.blockMap[block.GetHeight()] = block
}

func (b *BlockPool) Delete(height uint64) {
	b.mux.Lock()
	defer b.mux.Unlock()

	delete(b.blockMap, height)
}

func (b *BlockPool) GetByHeight(height uint64) blockchain.DefaultBlock {
	b.mux.Lock()
	defer b.mux.Unlock()

	if block, ok := b.blockMap[height]; ok {
		return block
	}

	return blockchain.DefaultBlock{}
}

func (b *BlockPool) GetSortedKeys() []blockchain.BlockHeight {
	keys := make([]blockchain.BlockHeight, 0)
	for h := range b.blockMap {
		keys = append(keys, h)
	}

	sortkeys.Uint64s(keys)

	return keys
}

func (b *BlockPool) Size() int {
	return len(b.blockMap)
}

