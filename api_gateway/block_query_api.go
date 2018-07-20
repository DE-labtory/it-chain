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

	"errors"
	"log"
	"sync"

	"github.com/it-chain/leveldb-wrapper"
)

var ErrGetCommitedBlock = errors.New("Error in getting commited block")
var ErrAddCommitingBlock = errors.New("Error in add block which is going to be commited")
var ErrNewBlockStorage = errors.New("Error in construct block storage")

type BlockQueryApi struct {
	blockPoolRepository     BlockPoolRepository
	commitedBlockRepository CommitedBlockRepository
}

func NewBlockQueryApi(blockPoolRepo BlockPoolRepository, commitedBlockRepo CommitedBlockRepository) *BlockQueryApi {
	return &BlockQueryApi{
		blockPoolRepository:     blockPoolRepo,
		commitedBlockRepository: commitedBlockRepo,
	}
}

// TODO:
func (q *BlockQueryApi) GetStagedBlockByHeight(height blockchain.BlockHeight) (blockchain.Block, error) {
	return nil, nil
}

// TODO:
func (q *BlockQueryApi) GetStagedBlockById(blockId string) (blockchain.Block, error) {
	return nil, nil
}
func (q *BlockQueryApi) GetLastCommitedBlock() (blockchain.Block, error) {
	defaultBlock, err := q.commitedBlockRepository.GetLastBlock()
	return &defaultBlock, err
}
func (q *BlockQueryApi) GetCommitedBlockByHeight(height blockchain.BlockHeight) (blockchain.Block, error) {
	defaultBlock, err := q.commitedBlockRepository.GetBlockByHeight(height)
	return &defaultBlock, err
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

func NewCommitedBlockRepositoryImpl(dbPath string) (*CommitedBlockRepositoryImpl, error) {
	validator := new(blockchain.DefaultValidator)
	db := leveldbwrapper.CreateNewDB(dbPath)
	opts := map[string]interface{}{}

	blockStorage, err := yggdrasill.NewBlockStorage(db, validator, opts)
	if err != nil {
		return nil, ErrNewBlockStorage
	}

	return &CommitedBlockRepositoryImpl{
		mux:                 &sync.RWMutex{},
		BlockStorageManager: blockStorage,
	}, nil
}

func (cbr *CommitedBlockRepositoryImpl) Save(block blockchain.DefaultBlock) error {
	cbr.mux.Lock()
	defer cbr.mux.Unlock()

	err := cbr.BlockStorageManager.AddBlock(&block)
	if err != nil {
		log.Fatal(err)
		return ErrAddCommitingBlock
	}
	return nil
}

func (cbr *CommitedBlockRepositoryImpl) GetLastBlock() (blockchain.DefaultBlock, error) {
	cbr.mux.Lock()
	defer cbr.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := cbr.BlockStorageManager.GetLastBlock(block)
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetCommitedBlock
	}

	return *block, nil
}
func (cbr *CommitedBlockRepositoryImpl) GetBlockByHeight(height uint64) (blockchain.DefaultBlock, error) {
	cbr.mux.Lock()
	defer cbr.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := cbr.BlockStorageManager.GetBlockByHeight(block, height)
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetCommitedBlock
	}

	return *block, nil
}
