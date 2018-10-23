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
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	leveldbwrapper "github.com/it-chain/leveldb-wrapper"
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

const (
	blockSealDB   = "block_seal"
	blockHeightDB = "block_height"
	utilDB        = "util"
	lastBlockKey  = "last_block"
)

type BlockQueryApi struct {
	blockRepository BlockRepository
}

func NewBlockQueryApi(blockRepo BlockRepository) *BlockQueryApi {
	return &BlockQueryApi{
		blockRepository: blockRepo,
	}
}

func (q BlockQueryApi) GetLastCommittedBlock() (BlockForQuery, error) {
	return q.blockRepository.FindLastBlock()
}

func (q BlockQueryApi) GetCommittedBlockByHeight(height blockchain.BlockHeight) (BlockForQuery, error) {
	return q.blockRepository.FindBlockByHeight(height)
}

func (q BlockQueryApi) GetCommittedBlockBySeal(seal []byte) (BlockForQuery, error) {
	return q.blockRepository.FindBlockBySeal(seal)
}

type BlockRepository interface {
	Save(block BlockForQuery) error
	FindLastBlock() (BlockForQuery, error)
	FindBlockByHeight(height uint64) (BlockForQuery, error)
	FindBlockBySeal(seal []byte) (BlockForQuery, error)
	FindAllBlock() ([]BlockForQuery, error)
	Close()
}

type BlockRepositoryImpl struct {
	mux        *sync.RWMutex
	DBProvider *DBProvider
}

func NewBlockRepositoryImpl(dbPath string) *BlockRepositoryImpl {

	db := leveldbwrapper.CreateNewDB(dbPath)
	dbProvider := CreateNewDBProvider(db)
	return &BlockRepositoryImpl{
		mux:        &sync.RWMutex{},
		DBProvider: dbProvider,
	}
}

func (r *BlockRepositoryImpl) Save(block BlockForQuery) error {
	r.mux.Lock()
	defer r.mux.Unlock()

	if block.Seal == nil {
		return ErrEmptyBlock
	}

	b, err := common.Serialize(block)
	if err != nil {
		return err
	}

	utilDB := r.DBProvider.GetDBHandle(utilDB)
	blockSealDB := r.DBProvider.GetDBHandle(blockSealDB)
	blockHeightDB := r.DBProvider.GetDBHandle(blockHeightDB)

	if err = blockSealDB.Put(block.Seal, b, true); err != nil {
		return err
	}

	if err = blockHeightDB.Put([]byte(fmt.Sprint(block.Height)), b, true); err != nil {
		return err
	}

	if err = utilDB.Put([]byte(lastBlockKey), b, true); err != nil {
		return err
	}

	return nil
}

func (r *BlockRepositoryImpl) FindLastBlock() (BlockForQuery, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	utilDB := r.DBProvider.GetDBHandle(utilDB)

	b, err := utilDB.Get([]byte(lastBlockKey))
	if err != nil {
		return BlockForQuery{}, err
	}

	if b == nil {
		return BlockForQuery{}, errors.New("Repository Is Empty")
	}

	block := &BlockForQuery{}

	if err = common.Deserialize(b, block); err != nil {
		return BlockForQuery{}, err
	}

	return *block, nil
}
func (r *BlockRepositoryImpl) FindBlockByHeight(height uint64) (BlockForQuery, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	blockHeightDB := r.DBProvider.GetDBHandle(blockHeightDB)

	b, err := blockHeightDB.Get([]byte(fmt.Sprint(height)))
	if err != nil {
		return BlockForQuery{}, err
	}

	block := &BlockForQuery{}

	if err = common.Deserialize(b, block); err != nil {
		return BlockForQuery{}, err
	}

	return *block, nil
}

func (r *BlockRepositoryImpl) FindBlockBySeal(seal []byte) (BlockForQuery, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	blockSealDB := r.DBProvider.GetDBHandle(blockSealDB)

	b, err := blockSealDB.Get(seal)
	if err != nil {
		return BlockForQuery{}, err
	}

	block := &BlockForQuery{}

	if err = common.Deserialize(b, block); err != nil {
		return BlockForQuery{}, err
	}

	return *block, nil

}

func (r *BlockRepositoryImpl) FindAllBlock() ([]BlockForQuery, error) {
	r.mux.Lock()
	defer r.mux.Unlock()

	blockHeightDB := r.DBProvider.GetDBHandle(blockHeightDB)

	iter := blockHeightDB.GetIteratorWithPrefix()

	blocks := []BlockForQuery{}

	for iter.Next() {
		val := iter.Value()
		block := &BlockForQuery{}
		err := common.Deserialize(val, block)
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, *block)
	}

	return blocks, nil
}

func (r *BlockRepositoryImpl) Close() {
	r.DBProvider.Close()
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

	block, err := createBlockForQuery(
		event.Seal,
		event.PrevSeal,
		event.Height,
		event.TxList,
		event.Timestamp,
		event.Creator,
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

func createBlockForQuery(Seal []byte, PrevSeal []byte, Height uint64, TxList []event.Tx, Timestamp time.Time, Creator string) (BlockForQuery, error) {
	txList, err := deserializeTxListType(TxList)
	if err != nil {
		return BlockForQuery{}, err
	}

	return BlockForQuery{
		Seal:      Seal,
		PrevSeal:  PrevSeal,
		Height:    Height,
		TxList:    txList,
		Timestamp: Timestamp,
		Creator:   Creator,
	}, nil
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

type BlockForQuery struct {
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []*blockchain.DefaultTransaction
	Timestamp time.Time
	Creator   string
}
