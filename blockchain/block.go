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

import (
	"encoding/json"
	"time"

	"bytes"

	e "github.com/it-chain/engine/common/event"
	"github.com/it-chain/midgard"
	ygg "github.com/it-chain/yggdrasill/common"

	"errors"
	"fmt"
	"reflect"
)

type Block = ygg.Block

type BlockHeight = uint64

type BlockState = string

const (
	Created   BlockState = "Created"
	Staged    BlockState = "Staged"
	Committed BlockState = "Committed"
)

type DefaultBlock struct {
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []*DefaultTransaction
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   []byte
	State     BlockState
}

// TODO: Write test case
func (block *DefaultBlock) SetSeal(seal []byte) {
	block.Seal = seal
}

// TODO: Write test case
func (block *DefaultBlock) SetPrevSeal(prevSeal []byte) {
	block.PrevSeal = prevSeal
}

// TODO: Write test case
func (block *DefaultBlock) SetHeight(height uint64) {
	block.Height = height
}

// TODO: Write test case
func (block *DefaultBlock) PutTx(transaction Transaction) error {
	convTx, ok := transaction.(*DefaultTransaction)
	if ok {
		if block.TxList == nil {
			block.TxList = make([]*DefaultTransaction, 0)
		}
		block.TxList = append(block.TxList, convTx)

		return nil
	}

	return ErrTransactionType
}

// TODO: Write test case
func (block *DefaultBlock) SetTxSeal(txSeal [][]byte) {
	block.TxSeal = txSeal
}

// TODO: Write test case
func (block *DefaultBlock) SetCreator(creator []byte) {
	block.Creator = creator
}

// TODO: Write test case
func (block *DefaultBlock) SetTimestamp(currentTime time.Time) {
	block.Timestamp = currentTime
}

// TODO: Write test case
func (block *DefaultBlock) GetSeal() []byte {
	return block.Seal
}

// TODO: Write test case
func (block *DefaultBlock) GetPrevSeal() []byte {
	return block.PrevSeal
}

// TODO: Write test case
func (block *DefaultBlock) GetHeight() uint64 {
	return block.Height
}

// TODO: Write test case
func (block *DefaultBlock) GetTxList() []Transaction {
	txList := make([]Transaction, 0)
	for _, tx := range block.TxList {
		txList = append(txList, tx)
	}
	return txList
}

// TODO: Write test case
func (block *DefaultBlock) GetTxSeal() [][]byte {
	return block.TxSeal
}

// TODO: Write test case
func (block *DefaultBlock) GetCreator() []byte {
	return block.Creator
}

// TODO: Write test case
func (block *DefaultBlock) GetTimestamp() time.Time {
	return block.Timestamp
}

// TODO: Write test case
func (block *DefaultBlock) Serialize() ([]byte, error) {
	data, err := json.Marshal(block)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// TODO: Write test case
func (block *DefaultBlock) Deserialize(serializedBlock []byte) error {
	if len(serializedBlock) == 0 {
		return ErrDecodingEmptyBlock
	}

	err := json.Unmarshal(serializedBlock, block)
	if err != nil {
		return err
	}

	return nil
}

// TODO: Write test case
func (block *DefaultBlock) IsReadyToPublish() bool {
	return block.Seal != nil
}

// TODO: Write test case
func (block *DefaultBlock) IsPrev(serializedPrevBlock []byte) bool {
	prevBlock := &DefaultBlock{}
	prevBlock.Deserialize(serializedPrevBlock)

	return bytes.Compare(prevBlock.GetSeal(), block.GetPrevSeal()) == 0
}

func (block *DefaultBlock) IsEmpty() bool {
	return reflect.DeepEqual(*block, DefaultBlock{})
}

func (block *DefaultBlock) GetBlockState() BlockState {
	return block.State
}

func (block *DefaultBlock) SetBlockState(state BlockState) {
	block.State = state
}

// interface of api gateway query api
type BlockQueryApi interface {
	GetBlockByHeight(blockHeight uint64) (Block, error)
	GetBlockBySeal(seal []byte) (Block, error)
	GetBlockByTxID(txid string) (Block, error)
	GetLastBlock() (Block, error)
}

func (block *DefaultBlock) On(event midgard.Event) error {

	switch v := event.(type) {

	case *(e.BlockCreated):
		TxList := ConvertToTransactionList(v.TxList)

		block.Seal = v.Seal
		block.PrevSeal = v.PrevSeal
		block.Height = v.Height
		block.TxList = TxList
		block.TxSeal = v.TxSeal
		block.Timestamp = v.Timestamp
		block.Creator = v.Creator
		block.State = v.State
		break
	case *(e.BlockStaged):
	case *(e.BlockCommitted):
		block.State = v.State
		break
	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

func NewEmptyBlock(prevSeal []byte, height uint64, creator []byte) *DefaultBlock {
	block := &DefaultBlock{}
	return block
}

func CreateBlock(block Block) error {
	// create BlockCreatedEvent
	// save it to event store
	return nil
}

func StageBlock(block Block) error {
	// create BlockStageEvent
	// save it to event store
	return nil
}
