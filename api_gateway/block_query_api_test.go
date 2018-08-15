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

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/api_gateway/test/mock"
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common/event"
	"github.com/stretchr/testify/assert"
)

func TestCommitedBlockRepositoryImpl(t *testing.T) {
	dbPath := "./.db"

	// when
	cbr, err := api_gateway.NewCommitedBlockRepositoryImpl(dbPath)

	// then
	assert.Equal(t, nil, err)

	defer func() {
		cbr.Close()
		os.RemoveAll(dbPath)
	}()

	// when
	block1 := mock.GetNewBlock([]byte("genesis"), 0)
	err = cbr.Save(*block1)
	// then
	assert.NoError(t, err)

	// when
	block2 := mock.GetNewBlock(block1.GetSeal(), 1)
	err = cbr.Save(*block2)
	// then
	assert.NoError(t, err)

	// when
	retrievedBlock, err := cbr.FindBlockByHeight(block2.Height)
	// then
	assert.NoError(t, err)
	assert.Equal(t, block2.GetSeal(), retrievedBlock.GetSeal())
	assert.Equal(t, block2.GetHeight(), retrievedBlock.GetHeight())
	assert.Equal(t, block2.GetPrevSeal(), retrievedBlock.GetPrevSeal())

	// when
	AllBlock, err4 := cbr.FindAllBlock()

	// then
	assert.NoError(t, err4)
	assert.Equal(t, 2, len(AllBlock))
	assert.Equal(t, block1.GetSeal(), AllBlock[0].GetSeal())
	assert.Equal(t, block2.GetSeal(), AllBlock[1].GetSeal())

}

func TestBlockEventListener_HandleBlockCommitedEvent(t *testing.T) {
	dbPath := "./.db"

	// given
	cbr, err := api_gateway.NewCommitedBlockRepositoryImpl(dbPath)

	assert.Equal(t, nil, err)

	defer func() {
		cbr.Close()
		os.RemoveAll(dbPath)
	}()

	lastCommittedBlock := mock.GetNewBlock([]byte("genesis"), 0)
	err = cbr.Save(*lastCommittedBlock)

	assert.NoError(t, err)

	eh := api_gateway.NewBlockEventListener(cbr)

	ConfirmedBlock := mock.GetNewBlock(lastCommittedBlock.GetSeal(), 1)
	ConfirmedBlock.State = blockchain.Committed

	txList, err := serializeTxListType(ConfirmedBlock.TxList)
	assert.NoError(t, err)

	blockCommittedEvent := event.BlockCommitted{
		Seal:      ConfirmedBlock.Seal,
		PrevSeal:  ConfirmedBlock.PrevSeal,
		Height:    ConfirmedBlock.Height,
		TxList:    txList,
		TxSeal:    ConfirmedBlock.TxSeal,
		Timestamp: ConfirmedBlock.Timestamp,
		Creator:   ConfirmedBlock.Creator,
		State:     ConfirmedBlock.State,
	}
	// when - Handle BlockCommited event
	err = eh.HandleBlockCommitedEvent(blockCommittedEvent)
	// then
	assert.NoError(t, err)

	// when - Test whether save target block to yggdrasill
	retrievedBlock, err := cbr.FindBlockByHeight(1)
	// then
	assert.NoError(t, err)
	assert.Equal(t, retrievedBlock.GetSeal(), ConfirmedBlock.GetSeal())
	assert.Equal(t, blockchain.Committed, retrievedBlock.State)

}

func serializeTxListType(txlist []*blockchain.DefaultTransaction) ([]event.Tx, error) {
	eventTxList := make([]event.Tx, 0)

	for _, tx := range txlist {
		eventTx, err := serializeTxType(tx)
		if err != nil {
			return eventTxList, err
		}
		eventTxList = append(eventTxList, eventTx)
	}

	return eventTxList, nil
}

func serializeTxType(tx *blockchain.DefaultTransaction) (event.Tx, error) {
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
