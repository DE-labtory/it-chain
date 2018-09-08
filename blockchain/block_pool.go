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

package blockchain

import "sync"

type BlockPool struct {
	blockMap map[uint64]Block
	mux      sync.Mutex
}

func NewBlockPool() *BlockPool {
	return &BlockPool{
		blockMap: make(map[uint64]Block),
		mux:      sync.Mutex{},
	}
}

func (b *BlockPool) Add(block Block) {
	b.mux.Lock()
	defer b.mux.Unlock()

	b.blockMap[block.GetHeight()] = block
}

func (b *BlockPool) Delete(height uint64) {
	b.mux.Lock()
	defer b.mux.Unlock()

	delete(b.blockMap, height)
}

func (b *BlockPool) Get(height uint64) Block {
	b.mux.Lock()
	defer b.mux.Unlock()

	if b.blockMap[height] == nil {
		return nil
	}

	return b.blockMap[height]
}
