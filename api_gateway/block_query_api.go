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

package api_gateway

import (
	"errors"
	"sync"
	"time"

	"github.com/DE-labtory/it-chain/blockchain"
	"github.com/DE-labtory/it-chain/common/event"
	leveldbwrapper "github.com/DE-labtory/leveldb-wrapper"
	"github.com/DE-labtory/yggdrasill"
)

var ErrGetCommittedBlock = errors.New("Error in getting commited block")
var ErrAddCommittingBlock = errors.New("Error in add block which is going to be commited")
var ErrNewBlockStorage = errors.New("Error in construct block storage")
var ErrNoCreatedBlock = errors.New("Error can not find created block")
var ErrNoStagedBlock = errors.New("Error can not find staged block")
var ErrInvalidStateBlock = errors.New("Error invalid state block")
var ErrFailRemoveBlock = errors.New("Error failed removing block")
var ErrIdEmpty = errors.New("Error that seal is empty string")
var ErrEmptyBlock = errors.New("Error empty block when getting block")

type BlockQueryApi struct {
	blockRepository BlockRepository
}

func NewBlockQueryApi(blockRepo BlockRepository) *BlockQueryApi {
	return &BlockQueryApi{
		blockRepository: blockRepo,
	}
}

func (q BlockQueryApi) GetLastCommittedBlock() (blockchain.DefaultBlock, error) {
	return q.blockRepository.FindLastBlock()
}

func (q BlockQueryApi) GetCommittedBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error) {
	return q.blockRepository.FindBlockByHeight(height)
}

func (q BlockQueryApi) GetCommittedBlockBySeal(seal []byte) (blockchain.DefaultBlock, error) {
	return q.blockRepository.FindBlockBySeal(seal)
}

type BlockRepository interface {
	Save(block blockchain.DefaultBlock) error
	FindLastBlock() (blockchain.DefaultBlock, error)
	FindBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
	FindBlockBySeal(seal []byte) (blockchain.DefaultBlock, error)
	FindAllBlock() ([]blockchain.DefaultBlock, error)
	Close()
}

type BlockRepositoryImpl struct {
	mux *sync.RWMutex
	yggdrasill.BlockStorageManager
}

func NewBlockRepositoryImpl(dbPath string) (*BlockRepositoryImpl, error) {
	validator := new(blockchain.DefaultValidator)
	db := leveldbwrapper.CreateNewDB(dbPath)
	opts := map[string]interface{}{}

	blockStorage, err := yggdrasill.NewBlockStorage(db, validator, opts)
	if err != nil {
		return nil, ErrNewBlockStorage
	}

	return &BlockRepositoryImpl{
		mux:                 &sync.RWMutex{},
		BlockStorageManager: blockStorage,
	}, nil
}

func (r *BlockRepositoryImpl) Save(block blockchain.DefaultBlock) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	err := r.BlockStorageManager.AddBlock(&block)
	if err != nil {
		return ErrAddCommittingBlock
	}

	return nil
}

func (r *BlockRepositoryImpl) FindLastBlock() (blockchain.DefaultBlock, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := r.BlockStorageManager.GetLastBlock(block)
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetCommittedBlock
	}

	return *block, nil
}
func (r *BlockRepositoryImpl) FindBlockByHeight(height uint64) (blockchain.DefaultBlock, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := r.BlockStorageManager.GetBlockByHeight(block, height)
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetCommittedBlock
	}

	return *block, nil
}

func (r *BlockRepositoryImpl) FindBlockBySeal(seal []byte) (blockchain.DefaultBlock, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := r.BlockStorageManager.GetBlockBySeal(block, seal)
	if err != nil {
		return blockchain.DefaultBlock{}, ErrGetCommittedBlock
	}

	return *block, nil
}

func (r *BlockRepositoryImpl) FindAllBlock() ([]blockchain.DefaultBlock, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	blocks := []blockchain.DefaultBlock{}

	// set
	lastBlock := &blockchain.DefaultBlock{}

	err := r.BlockStorageManager.GetLastBlock(lastBlock)

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

		err := r.BlockStorageManager.GetBlockByHeight(block, i)

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
	blockRepository BlockRepository
}

func NewBlockEventListener(blockRepo BlockRepository) *BlockEventListener {
	return &BlockEventListener{
		blockRepository: blockRepo,
	}
}

func (l BlockEventListener) HandleBlockCommittedEvent(event event.BlockCommitted) error {
	blockID := event.Seal
	if len(blockID) == 0 {
		return ErrIdEmpty
	}

	block, err := createDefaultBlock(
		event.Seal,
		event.PrevSeal,
		event.Height,
		event.TxList,
		event.TxSeal,
		event.Timestamp,
		event.Creator,
		event.State,
	)
	if err != nil {
		return err
	}

	err = l.blockRepository.Save(block)
	if err != nil {
		return err
	}

	return nil
}

func createDefaultBlock(Seal []byte, PrevSeal []byte, Height uint64, TxList []event.Tx, TxSeal [][]byte, Timestamp time.Time, Creator string, State string) (blockchain.DefaultBlock, error) {
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
