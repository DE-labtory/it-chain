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

	"errors"
	"sync"

	"time"

	"bytes"

	"encoding/hex"

	"encoding/json"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/leveldb-wrapper"
)

var ErrNoCreatedBlock = errors.New("Error can not find created block")
var ErrNoStagedBlock = errors.New("Error can not find staged block")
var ErrNoBlock = errors.New("Error can not find block")
var ErrInvalidStateBlock = errors.New("Error invalid state block")
var ErrFailRemoveBlock = errors.New("Error failed removing block")
var ErrFailBlockTypeCasting = errors.New("Error failed type casting block")
var ErrIdEmpty = errors.New("Error that seal is empty string")
var ErrEmptyBlock = errors.New("Error empty block when not expected")

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

type BlockPoolRepository interface {
	SaveCreatedBlock(block blockchain.DefaultBlock) error
	SaveStagedBlock(block blockchain.DefaultBlock) error
	FindCreatedBlockById(id string) (blockchain.DefaultBlock, error)
	FindStagedBlockById(id string) (blockchain.DefaultBlock, error)
	FindStagedBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
	FindFirstStagedBlock() (blockchain.DefaultBlock, error)
	FindBlockById(id string) (blockchain.DefaultBlock, error)
	RemoveBlockById(id string) error
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
		if isBlockSealEqualWith(b, block.GetSeal()) {
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
		if isBlockSealEqualWith(b, block.GetSeal()) {
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

func (r *BlockPoolRepositoryImpl) FindBlockById(id string) (blockchain.DefaultBlock, error) {
	for _, block := range r.Blocks {
		if isBlockIdEqualWith(block, id) {

			defaultBlock, ok := block.(*blockchain.DefaultBlock)
			if !ok {
				return blockchain.DefaultBlock{}, ErrFailBlockTypeCasting
			}

			return *defaultBlock, nil
		}
	}

	return blockchain.DefaultBlock{}, ErrNoBlock
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

func (r *BlockPoolRepositoryImpl) RemoveBlockById(id string) error {
	for i, b := range r.Blocks {
		if isBlockIdEqualWith(b, id) {
			r.Blocks = append(r.Blocks[:i], r.Blocks[i+1:]...)
			return nil
		}
	}
	return ErrFailRemoveBlock
}

func isBlockSealEqualWith(block blockchain.Block, seal []byte) bool {

	return bytes.Equal(block.GetSeal(), seal)
}

func isBlockIdEqualWith(block blockchain.Block, id string) bool {
	BlockID := hex.EncodeToString(block.GetSeal())

	return BlockID == id
}

type CommitedBlockRepository interface {
	Save(block blockchain.DefaultBlock) error
	FindBlockByHeight(height uint64) (blockchain.DefaultBlock, error)
	FindAllBlock() ([]blockchain.DefaultBlock, error)
}

type CommitedBlockRepositoryImpl struct {
	mux     *sync.RWMutex
	leveldb *leveldbwrapper.DB
}

func NewCommitedBlockRepositoryImpl(dbPath string) (*CommitedBlockRepositoryImpl, error) {
	db := leveldbwrapper.CreateNewDB(dbPath)
	db.Open()

	return &CommitedBlockRepositoryImpl{
		mux:     &sync.RWMutex{},
		leveldb: db,
	}, nil
}

func (cbr *CommitedBlockRepositoryImpl) Save(block blockchain.DefaultBlock) error {
	cbr.mux.Lock()
	defer cbr.mux.Unlock()

	if block.IsEmpty() {
		return ErrEmptyBlock
	}

	b, err := block.Serialize()

	if err != nil {
		return err
	}

	key := []byte(string(block.Height))

	if err = cbr.leveldb.Put(key, b, true); err != nil {
		return err
	}

	return nil
}

func (cbr *CommitedBlockRepositoryImpl) FindBlockByHeight(height uint64) (blockchain.DefaultBlock, error) {
	cbr.mux.Lock()
	defer cbr.mux.Unlock()

	key := []byte(string(height))

	b, err := cbr.leveldb.Get(key)

	if err != nil {
		return blockchain.DefaultBlock{}, err
	}

	block := blockchain.DefaultBlock{}

	err = json.Unmarshal(b, &block)

	if err != nil {
		return blockchain.DefaultBlock{}, err
	}

	return block, nil

}

func (cbr *CommitedBlockRepositoryImpl) FindAllBlock() ([]blockchain.DefaultBlock, error) {
	cbr.mux.Lock()
	defer cbr.mux.Unlock()

	iter := cbr.leveldb.GetIteratorWithPrefix([]byte(""))

	blocks := []blockchain.DefaultBlock{}

	for iter.Next() {

		val := iter.Value()

		block := &blockchain.DefaultBlock{}

		err := block.Deserialize(val)

		if err != nil {
			return nil, err
		}

		blocks = append(blocks, *block)
	}

	return blocks, nil
}

func (cbr *CommitedBlockRepositoryImpl) Close() {
	cbr.leveldb.Close()
}

type BlockEventListener struct {
	blockPoolRepository     BlockPoolRepository
	commitedBlockRepository CommitedBlockRepository
}

func NewBlockEventListener(blockPoolRepo BlockPoolRepository, commitedBlockRepo CommitedBlockRepository) BlockEventListener {
	return BlockEventListener{
		blockPoolRepository:     blockPoolRepo,
		commitedBlockRepository: commitedBlockRepo,
	}
}

func (l BlockEventListener) HandleBlockCreatedEvent(event event.BlockCreated) error {
	block, err := createDefaultBlock(event.Seal, event.PrevSeal, event.Height, event.TxList, event.TxSeal, event.Timestamp, event.Creator, event.State)
	if err != nil {
		return err
	}

	l.blockPoolRepository.SaveCreatedBlock(block)

	return nil
}

func (l BlockEventListener) HandleBlockStagedEvent(event event.BlockStaged) error {
	blockID := event.ID
	if blockID == "" {
		return ErrIdEmpty
	}

	block, err := l.blockPoolRepository.FindCreatedBlockById(blockID)
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

func (l BlockEventListener) HandleBlockCommitedEvent(event event.BlockCommitted) error {

	txList, err := deserializeTxListType(event.TxList)

	if err != nil {
		return err
	}

	block := blockchain.DefaultBlock{
		Seal:      event.Seal,
		PrevSeal:  event.PrevSeal,
		Height:    event.Height,
		TxList:    txList,
		TxSeal:    event.TxSeal,
		Timestamp: event.Timestamp,
		Creator:   event.Creator,
		State:     event.State,
	}

	err = l.commitedBlockRepository.Save(block)
	if err != nil {
		return err
	}

	return nil
}

func createDefaultBlock(Seal []byte, PrevSeal []byte, Height uint64, TxList []event.Tx, TxSeal [][]byte, Timestamp time.Time, Creator []byte, State string) (blockchain.DefaultBlock, error) {
	txList, err := deserializeTxListType(TxList)
	if err != nil {
		return blockchain.DefaultBlock{}, err
	}
	return blockchain.DefaultBlock{
		Seal:      Seal,
		PrevSeal:  PrevSeal,
		Height:    Height,
		TxList:    txList,
		TxSeal:    TxSeal,
		Timestamp: Timestamp,
		Creator:   Creator,
		State:     State,
	}, nil
}

func deserializeTxListType(txlist []event.Tx) ([]*blockchain.DefaultTransaction, error) {
	defaultTxList := make([]*blockchain.DefaultTransaction, 0)

	for _, tx := range txlist {
		defaultTx, err := deserializeTxType(tx)
		if err != nil {
			return defaultTxList, err
		}
		defaultTxList = append(defaultTxList, defaultTx)
	}

	return defaultTxList, nil
}

func deserializeTxType(tx event.Tx) (*blockchain.DefaultTransaction, error) {
	return &blockchain.DefaultTransaction{
		ID:        tx.ID,
		ICodeID:   tx.ICodeID,
		PeerID:    tx.PeerID,
		Timestamp: tx.TimeStamp,
		Jsonrpc:   tx.Jsonrpc,
		Function:  tx.Function,
		Args:      tx.Args,
		Signature: tx.Signature,
	}, nil
}
