/*
 * Copyright 2018 DE-labtory
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

package repo

import (
	"sync"

	"github.com/DE-labtory/engine/blockchain"
	"github.com/DE-labtory/leveldb-wrapper"
	"github.com/DE-labtory/yggdrasill"
)

type BlockRepository struct {
	mux *sync.RWMutex
	yggdrasill.BlockStorageManager
}

func NewBlockRepository(dbPath string) (*BlockRepository, error) {
	validator := new(blockchain.DefaultValidator)
	db := leveldbwrapper.CreateNewDB(dbPath)
	opts := map[string]interface{}{}

	blockStorage, err := yggdrasill.NewBlockStorage(db, validator, opts)
	if err != nil {
		return nil, ErrNewBlockStorage
	}

	return &BlockRepository{
		mux:                 &sync.RWMutex{},
		BlockStorageManager: blockStorage,
	}, nil
}

func (br *BlockRepository) Save(block blockchain.DefaultBlock) error {
	br.mux.Lock()
	defer br.mux.Unlock()
	err := br.BlockStorageManager.AddBlock(&block)
	if err != nil {
		return ErrAddBlock
	}

	return nil
}

func (br *BlockRepository) FindLast() (blockchain.DefaultBlock, error) {
	br.mux.Lock()
	defer br.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := br.BlockStorageManager.GetLastBlock(block)
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetBlock
	}

	return *block, nil
}
func (br *BlockRepository) FindByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
	br.mux.Lock()
	defer br.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := br.BlockStorageManager.GetBlockByHeight(block, height)
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetBlock
	}

	return *block, nil
}

func (br *BlockRepository) FindBySeal(seal []byte) (blockchain.DefaultBlock, error) {
	br.mux.Lock()
	defer br.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := br.BlockStorageManager.GetBlockBySeal(block, seal)
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetBlock
	}

	return *block, nil
}

func (br *BlockRepository) FindAll() ([]blockchain.DefaultBlock, error) {
	br.mux.Lock()
	defer br.mux.Unlock()

	blocks := []blockchain.DefaultBlock{}

	// set
	lastBlock := &blockchain.DefaultBlock{}

	err := br.BlockStorageManager.GetLastBlock(lastBlock)

	if err != nil {
		return nil, err
	}

	// check empty
	if lastBlock.IsEmpty() {
		return blocks, nil
	}

	lastHeight := lastBlock.GetHeight()

	// get blocks
	for i := uint64(0); i <= lastHeight; i++ {

		block := &blockchain.DefaultBlock{}

		err := br.BlockStorageManager.GetBlockByHeight(block, i)

		if err != nil {
			return nil, err
		}

		if block.IsEmpty() {
			return nil, ErrEmptyBlock
		}

		blocks = append(blocks, *block)
	}

	return blocks, nil
}
