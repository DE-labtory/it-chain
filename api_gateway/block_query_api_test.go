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

package api_gateway_test

import (
	"os"
	"testing"

	"encoding/hex"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/api_gateway/test/mock"
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common/event"
	"github.com/stretchr/testify/assert"
)

func TestBlockQueryApi_FindLastCommitedBlock(t *testing.T) {
	dbPath := "./.db"

	// when
	cbr, err := api_gateway.NewBlockRepositoryImpl(dbPath)
	// then
	assert.Equal(t, nil, err)

	defer func() {
		cbr.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	err = cbr.AddBlock(block1)
	// then
	assert.NoError(t, err)

	// when
	block2 := mock.GetNewBlock(block1.GetSeal(), 1)
	err = cbr.AddBlock(block2)
	// then
	assert.NoError(t, err)

	blockQueryApi := api_gateway.NewBlockQueryApi(cbr)

	// when
	block3, err := blockQueryApi.GetLastCommittedBlock()
	// then
	assert.NoError(t, err)
	assert.Equal(t, block2.GetSeal(), block3.GetSeal())
	assert.Equal(t, block2.GetHeight(), block3.GetHeight())
	assert.Equal(t, block2.GetPrevSeal(), block3.GetPrevSeal())
}

func TestBlockQueryApi_FindCommitedBlockByHeight(t *testing.T) {
	dbPath := "./.db"

	// when
	cbr, err := api_gateway.NewBlockRepositoryImpl(dbPath)
	// then
	assert.Equal(t, nil, err)

	defer func() {
		cbr.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	err = cbr.AddBlock(block1)
	// then
	assert.NoError(t, err)

	// when
	block2 := mock.GetNewBlock(block1.GetSeal(), 1)
	err = cbr.AddBlock(block2)
	// then
	assert.NoError(t, err)

	blockQueryApi := api_gateway.NewBlockQueryApi(cbr)

	// when
	block3, err := blockQueryApi.GetCommittedBlockByHeight(blockchain.BlockHeight(1))
	// then
	assert.NoError(t, err)
	assert.Equal(t, block2.GetSeal(), block3.GetSeal())
	assert.Equal(t, block2.GetHeight(), block3.GetHeight())
	assert.Equal(t, block2.GetPrevSeal(), block3.GetPrevSeal())
}

func TestCommitedBlockRepositoryImpl(t *testing.T) {
	dbPath := "./.db"

	// when
	cbr, err := api_gateway.NewBlockRepositoryImpl(dbPath)

	// then
	assert.Equal(t, nil, err)

	defer func() {
		cbr.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	err = cbr.AddBlock(block1)
	// then
	assert.NoError(t, err)

	// when
	block2 := mock.GetNewBlock(block1.GetSeal(), 1)
	err = cbr.AddBlock(block2)
	// then
	assert.NoError(t, err)

	blockQueryApi := api_gateway.NewBlockQueryApi(cbr)

	// when
	block3, err := blockQueryApi.GetLastCommittedBlock()
	// then
	assert.NoError(t, err)
	assert.Equal(t, block2.GetSeal(), block3.GetSeal())
	assert.Equal(t, block2.GetHeight(), block3.GetHeight())
	assert.Equal(t, block2.GetPrevSeal(), block3.GetPrevSeal())

	// when
	AllBlock, err4 := cbr.FindAllBlock()

	// then
	assert.NoError(t, err4)
	assert.Equal(t, 2, len(AllBlock))

}

func TestBlockEventListener_HandleBlockCommitedEvent(t *testing.T) {
	dbPath := "./.db"

	// when
	cbr, err := api_gateway.NewBlockRepositoryImpl(dbPath)
	// then
	assert.Equal(t, nil, err)

	defer func() {
		cbr.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	err = cbr.AddBlock(block1)
	// then
	assert.NoError(t, err)

	eh := api_gateway.NewBlockEventListener(cbr)

	// when
	block2 := mock.GetNewBlock(block1.GetSeal(), 1)
	// when
	block2ID := hex.EncodeToString(block2.GetSeal())
	txList, _ := convertToTxList(block2.TxList)

	event1 := event.BlockCommitted{
		BlockId:   block2ID,
		Seal:      block2.Seal,
		PrevSeal:  block2.PrevSeal,
		Height:    block2.Height,
		TxList:    txList,
		TxSeal:    block2.TxSeal,
		Timestamp: block2.Timestamp,
		Creator:   block2.Creator,
		State:     blockchain.Committed,
	}
	// when - Handle BlockCommited event
	err1 := eh.HandleBlockCommittedEvent(event1)
	// then
	assert.NoError(t, err1)

	// when - Test whether save target block to yggdrasill
	block3, err2 := cbr.FindBlockByHeight(1)
	// then
	assert.NoError(t, err2)
	assert.Equal(t, block3.Seal, block2.GetSeal())
	assert.Equal(t, blockchain.Committed, block3.State)
}

func convertToTxList(txlist []*blockchain.DefaultTransaction) ([]event.Tx, error) {
	defaultTxList := make([]event.Tx, 0)

	for _, tx := range txlist {
		defaultTx, err := convertToTx(tx)
		if err != nil {
			return defaultTxList, err
		}
		defaultTxList = append(defaultTxList, defaultTx)
	}

	return defaultTxList, nil
}

func convertToTx(tx *blockchain.DefaultTransaction) (event.Tx, error) {
	return event.Tx{
		ID:        tx.ID,
		ICodeID:   tx.ICodeID,
		PeerID:    tx.PeerID,
		TimeStamp: tx.Timestamp,
		Jsonrpc:   tx.Jsonrpc,
		Function:  tx.Function,
		Args:      tx.Args,
		Signature: tx.Signature,
	}, nil
}
