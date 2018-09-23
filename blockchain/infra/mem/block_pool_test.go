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

package mem_test

import (
	"testing"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/mem"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/stretchr/testify/assert"
)

func TestBlockPool(t *testing.T) {
	blockPool := mem.NewBlockPool()

	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	block2 := mock.GetNewBlock(block1.GetSeal(), 1)
	block3 := mock.GetNewBlock(block2.GetSeal(), 2)

	blockPool.Add(*block1)
	blockPool.Add(*block2)
	blockPool.Add(*block3)
	assert.Equal(t, blockPool.Size(), 3)

	b := blockPool.GetByHeight(1)
	assert.Equal(t, b.GetHeight(), uint64(1))
	assert.Equal(t, b.GetSeal(), (*block2).GetSeal())

	b2 := blockPool.GetByHeight(123123)
	assert.Equal(t, b2, blockchain.DefaultBlock{})

	blockPool.Delete(1)
	assert.Equal(t, blockPool.Size(), 2)

	b3 := blockPool.GetByHeight(1)
	assert.Equal(t, b3, blockchain.DefaultBlock{})
}

func TestBlockPool_GetSortedKeys(t *testing.T) {
	blockPool := mem.NewBlockPool()

	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	block2 := mock.GetNewBlock(block1.GetSeal(), 3120)
	block3 := mock.GetNewBlock(block2.GetSeal(), 29)
	block4 := mock.GetNewBlock(block2.GetSeal(), 694)

	blockPool.Add(*block1)
	blockPool.Add(*block2)
	blockPool.Add(*block3)
	blockPool.Add(*block4)

	sortedKeys := blockPool.GetSortedKeys()
	assert.Equal(t, sortedKeys, []uint64{uint64(0), uint64(29), uint64(694), uint64(3120)})
}
