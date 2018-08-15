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

	"encoding/json"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/leveldb-wrapper"
)

var ErrEmptyBlock = errors.New("Error empty block when not expected")

type BlockQueryApi struct {
	commitedBlockRepository CommitedBlockRepository
}

func NewBlockQueryApi(commitedBlockRepo CommitedBlockRepository) BlockQueryApi {
	return BlockQueryApi{
		commitedBlockRepository: commitedBlockRepo,
	}
}

type CommitedBlockRepository interface {
	Save(block blockchain.DefaultBlock) error
	FindBlockByHeight(height uint64) (*blockchain.DefaultBlock, error)
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

func (cbr *CommitedBlockRepositoryImpl) FindBlockByHeight(height uint64) (*blockchain.DefaultBlock, error) {
	cbr.mux.Lock()
	defer cbr.mux.Unlock()

	key := []byte(string(height))

	b, err := cbr.leveldb.Get(key)

	if err != nil {
		return nil, err
	}

	block := &blockchain.DefaultBlock{}

	err = json.Unmarshal(b, block)

	if err != nil {
		return nil, err
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
	commitedBlockRepository CommitedBlockRepository
}

func NewBlockEventListener(commitedBlockRepo CommitedBlockRepository) BlockEventListener {
	return BlockEventListener{
		commitedBlockRepository: commitedBlockRepo,
	}
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
