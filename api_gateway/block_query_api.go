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

	"time"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/leveldb-wrapper"
)

var ErrGetCommitedBlock = errors.New("Error in getting commited block")
var ErrAddCommitingBlock = errors.New("Error in add block which is going to be commited")
var ErrNewBlockStorage = errors.New("Error in construct block storage")
var ErrNoCreatedBlock = errors.New("Error can not find created block")
var ErrNoStagedBlock = errors.New("Error can not find staged block")
var ErrInvalidStateBlock = errors.New("Error invalid state block")
var ErrFailRemoveBlock = errors.New("Error failed removing block")
var ErrFailBlockTypeCasting = errors.New("Error failed type casting block")
var ErrSealEmpty = errors.New("Error that seal is empty string")
var ErrEmptyBlock = errors.New("Error empty block when getting block")
var ErrCheckEmpty = errors.New("Error when checking repo empty")

type BlockQueryApi struct {
	blockPoolRepository     BlockPoolRepository
	commitedBlockRepository CommitedBlockRepository
}

func NewBlockQueryApi(blockPoolRepo BlockPoolRepository, commitedBlockRepo CommitedBlockRepository) BlockQueryApi {
	return BlockQueryApi{
		blockPoolRepository:     blockPoolRepo,
		commitedBlockRepository: commitedBlockRepo,
	}
}

func (q BlockQueryApi) GetStagedBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
	return q.blockPoolRepository.FindStagedBlockByHeight(height)
}
func (q BlockQueryApi) GetStagedBlockById(blockId string) (blockchain.DefaultBlock, error) {
	return q.blockPoolRepository.FindStagedBlockById(blockId)
}

func (q BlockQueryApi) GetLastCommitedBlock() (blockchain.DefaultBlock, error) {
	return q.commitedBlockRepository.FindLast()
}
func (q BlockQueryApi) GetCommitedBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
	return q.commitedBlockRepository.FindByHeight(height)
}

type BlockPoolRepository interface {
	SaveCreatedBlock(block blockchain.DefaultBlock) error
	SaveStagedBlock(block blockchain.DefaultBlock) error
	FindCreatedBlockById(id string) (blockchain.DefaultBlock, error)
	FindStagedBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
	FindStagedBlockById(id string) (blockchain.DefaultBlock, error)
	FindFirstStagedBlock() (blockchain.DefaultBlock, error)
	RemoveById(id string) error
}

type BlockPoolRepositoryImpl struct {
	Blocks []blockchain.Block
}

func NewBlockPoolRepository() *BlockPoolRepositoryImpl {
	return &BlockPoolRepositoryImpl{
		Blocks: make([]blockchain.Block, 0),
	}
}

func (r *BlockPoolRepositoryImpl) SaveCreatedBlock(block blockchain.DefaultBlock) error {
	if block.State != blockchain.Created {
		return ErrInvalidStateBlock
	}

	for i, b := range r.Blocks {
		if isBlockIdEqualWith(b, string(block.GetSeal())) {
			r.Blocks[i] = &block
			return nil
		}
	}
	r.Blocks = append(r.Blocks, &block)
	return nil
}

func (r *BlockPoolRepositoryImpl) SaveStagedBlock(block blockchain.DefaultBlock) error {
	if block.State != blockchain.Staged {
		return ErrInvalidStateBlock
	}

	for i, b := range r.Blocks {
		if isBlockIdEqualWith(b, string(block.GetSeal())) {
			r.Blocks[i] = &block
			return nil
		}
	}
	r.Blocks = append(r.Blocks, &block)
	return nil
}

func (r *BlockPoolRepositoryImpl) FindCreatedBlockById(id string) (blockchain.DefaultBlock, error) {
	for _, block := range r.Blocks {
		if isBlockIdEqualWith(block, id) {
			defaultBlock, ok := block.(*blockchain.DefaultBlock)
			if !ok {
				return blockchain.DefaultBlock{}, ErrFailBlockTypeCasting
			}

			if defaultBlock.State == blockchain.Created {
				return *defaultBlock, nil
			}
		}
	}
	return blockchain.DefaultBlock{}, ErrNoCreatedBlock
}

func (r *BlockPoolRepositoryImpl) FindStagedBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
	for _, block := range r.Blocks {
		if block.GetHeight() == height {
			defaultBlock, ok := block.(*blockchain.DefaultBlock)
			if !ok {
				return blockchain.DefaultBlock{}, ErrFailBlockTypeCasting
			}

			if defaultBlock.State == blockchain.Staged {
				return *defaultBlock, nil
			}
		}
	}
	return blockchain.DefaultBlock{}, ErrNoStagedBlock
}

func (r *BlockPoolRepositoryImpl) FindStagedBlockById(id string) (blockchain.DefaultBlock, error) {
	for _, block := range r.Blocks {
		if isBlockIdEqualWith(block, id) {
			defaultBlock, ok := block.(*blockchain.DefaultBlock)
			if !ok {
				return blockchain.DefaultBlock{}, ErrFailBlockTypeCasting
			}

			if defaultBlock.State == blockchain.Staged {
				return *defaultBlock, nil
			}
		}
	}
	return blockchain.DefaultBlock{}, ErrNoStagedBlock
}

func (r *BlockPoolRepositoryImpl) FindFirstStagedBlock() (blockchain.DefaultBlock, error) {
	if len(r.Blocks) == 0 {
		return blockchain.DefaultBlock{}, ErrNoStagedBlock
	}

	target := blockchain.DefaultBlock{}

	for _, block := range r.Blocks {
		defaultBlock, ok := block.(*blockchain.DefaultBlock)
		if !ok {
			return blockchain.DefaultBlock{}, ErrFailBlockTypeCasting
		}

		if stagedBlockWithSmallerHeight(target, *defaultBlock) {
			target = *defaultBlock
		}
	}

	if target.IsEmpty() {
		return target, ErrNoStagedBlock
	}

	return target, nil
}

func stagedBlockWithSmallerHeight(base blockchain.DefaultBlock, comparator blockchain.DefaultBlock) bool {
	return comparator.State == blockchain.Staged && (base.Height > comparator.Height || base.IsEmpty())
}

func (r *BlockPoolRepositoryImpl) RemoveById(id string) error {
	for i, b := range r.Blocks {
		if isBlockIdEqualWith(b, id) {
			r.Blocks = append(r.Blocks[:i], r.Blocks[i+1:]...)
			return nil
		}
	}
	return ErrFailRemoveBlock
}

func isBlockIdEqualWith(block blockchain.Block, id string) bool {
	return string(block.GetSeal()) == id
}

type CommitedBlockRepository interface {
	Save(block blockchain.DefaultBlock) error
	FindLast() (blockchain.DefaultBlock, error)
	FindByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
	FindAll() ([]blockchain.DefaultBlock, error)
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

func (cbr *CommitedBlockRepositoryImpl) FindLast() (blockchain.DefaultBlock, error) {
	cbr.mux.Lock()
	defer cbr.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := cbr.BlockStorageManager.GetLastBlock(block)
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetCommitedBlock
	}

	return *block, nil
}
func (cbr *CommitedBlockRepositoryImpl) FindByHeight(height uint64) (blockchain.DefaultBlock, error) {
	cbr.mux.Lock()
	defer cbr.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := cbr.BlockStorageManager.GetBlockByHeight(block, height)
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetCommitedBlock
	}

	return *block, nil
}

func (cbr *CommitedBlockRepositoryImpl) FindAll() ([]blockchain.DefaultBlock, error) {
	cbr.mux.Lock()
	defer cbr.mux.Unlock()

	blocks := []blockchain.DefaultBlock{}

	// set
	lastBlock := &blockchain.DefaultBlock{}

	err := cbr.BlockStorageManager.GetLastBlock(lastBlock)

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

		err := cbr.BlockStorageManager.GetBlockByHeight(block, i)

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

type BlockEventListener struct {
	blockPoolRepository     BlockPoolRepository
	commitedBlockRepository CommitedBlockRepository
}

func NewBlockEventListener(blockPoolRepo BlockPoolRepository, commitedBlockRepo CommitedBlockRepository) *BlockEventListener {
	return &BlockEventListener{
		blockPoolRepository:     blockPoolRepo,
		commitedBlockRepository: commitedBlockRepo,
	}
}

func (l *BlockEventListener) HandleBlockCreatedEvent(event event.BlockCreated) error {
	block, err := createDefaultBlock(event.Seal, event.PrevSeal, event.Height, event.TxList, event.Timestamp, event.Creator, event.State)
	if err != nil {
		return err
	}

	l.blockPoolRepository.SaveCreatedBlock(block)
	return nil
}

func createDefaultBlock(Seal []byte, PrevSeal []byte, Height uint64, TxList []event.Tx, Timestamp time.Time, Creator []byte, State string) (blockchain.DefaultBlock, error) {
	txList, err := getBackTxType(TxList)
	if err != nil {
		return blockchain.DefaultBlock{}, err
	}
	return blockchain.DefaultBlock{
		Seal:      Seal,
		PrevSeal:  PrevSeal,
		Height:    Height,
		TxList:    txList,
		Timestamp: Timestamp,
		Creator:   Creator,
		State:     State,
	}, nil
}

func deserializeTxList(txList []byte) ([]*blockchain.DefaultTransaction, error) {
	DefaultTxList := []*blockchain.DefaultTransaction{}
	err := common.Deserialize(txList, &DefaultTxList)
	// get blocks
	if err != nil {
		return nil, err
	}
	return DefaultTxList, nil
}

func getBackTxType(txlist []event.Tx) ([]*blockchain.DefaultTransaction, error) {
	defaultTxList := make([]*blockchain.DefaultTransaction, 0)

	for _, tx := range txlist {
		defaultTx, err := convertToDefaultTransaction(tx)
		if err != nil {
			return defaultTxList, err
		}
		defaultTxList = append(defaultTxList, defaultTx)
	}

	return defaultTxList, nil
}

func convertToDefaultTransaction(tx event.Tx) (*blockchain.DefaultTransaction, error) {
	return &blockchain.DefaultTransaction{
		ID:        tx.ID,
		Status:    blockchain.Status(tx.Status),
		PeerID:    tx.PeerID,
		ICodeID:   tx.ICodeID,
		Timestamp: tx.TimeStamp,
		TxData: blockchain.TxData{
			Jsonrpc: tx.Jsonrpc,
			Method:  blockchain.TxDataType(tx.Method),
			Params: blockchain.Params{
				Function: tx.Function,
				Args:     tx.Args,
			},
			ID: tx.ID,
		},
		Signature: tx.Signature,
	}, nil
}

func (l *BlockEventListener) HandleBlockStagedEvent(event event.BlockStaged) error {
	seal := event.ID
	if seal == "" {
		return ErrSealEmpty
	}

	block, err := l.blockPoolRepository.FindCreatedBlockById(seal)
	if err != nil {
		return err
	}

	block.State = event.State

	err = l.blockPoolRepository.SaveStagedBlock(block)
	if err != nil {
		return err
	}

	return nil
}

func (l *BlockEventListener) HandleBlockCommitedEvent(event event.BlockCommitted) error {
	seal := event.ID
	if seal == "" {
		return ErrSealEmpty
	}

	block, err := l.blockPoolRepository.FindStagedBlockById(seal)
	if err != nil {
		return err
	}

	block.State = event.State

	err = l.commitedBlockRepository.Save(block)
	if err != nil {
		return err
	}

	err = l.blockPoolRepository.RemoveById(seal)
	if err != nil {
		return err
	}

	return nil
}
