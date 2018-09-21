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

package blockchain_test

import (
	"testing"

	"time"

	"github.com/it-chain/engine/blockchain"
	"github.com/stretchr/testify/assert"
)

func TestDefaultBlock_BasicFunctions(t *testing.T) {

	//given
	Creator := "junksound"
	TimeStamp := time.Now().Round(0)
	State := blockchain.Created

	TxList := []*blockchain.DefaultTransaction{
		{
			ID:        "tx01",
			ICodeID:   "Icode01",
			PeerID:    "Peer01",
			Timestamp: TimeStamp,
			Jsonrpc:   "jsonrpc01",
			Function:  "function01",
			Args:      nil,
			Signature: []byte("Signature"),
		},
		{
			ID:        "tx02",
			ICodeID:   "Icode02",
			PeerID:    "Peer02",
			Timestamp: TimeStamp,
			Jsonrpc:   "jsonrpc02",
			Function:  "function02",
			Args:      nil,
			Signature: []byte("Signature"),
		},
	}

	Block := blockchain.DefaultBlock{}
	validator := blockchain.DefaultValidator{}

	//when
	for _, tx := range TxList {
		Block.PutTx(tx)
	}

	convertedTxList := convertTxType(TxList)

	txSeal, err := validator.BuildTxSeal(convertedTxList)
	assert.NoError(t, err)

	Block.SetTxSeal(txSeal)
	Block.SetCreator(Creator)
	Block.SetTimestamp(TimeStamp)
	Block.SetState(State)
	Block.SetSeal([]byte("seal"))
	Block.SetPrevSeal([]byte("prevSeal"))
	Block.SetHeight(1)
	assert.Equal(t, uint64(1), Block.GetHeight())

	//then
	assert.Equal(t, convertedTxList, Block.GetTxList())
	assert.Equal(t, txSeal, Block.GetTxSeal())
	assert.Equal(t, Creator, Block.GetCreator())
	assert.Equal(t, TimeStamp, Block.GetTimestamp())
	assert.Equal(t, State, Block.GetState())
	assert.Equal(t, []byte("seal"), Block.GetSeal())
	assert.Equal(t, []byte("prevSeal"), Block.GetPrevSeal())
	assert.Equal(t, uint64(1), Block.GetHeight())
	assert.Equal(t, false, Block.IsEmpty())
}

func convertTxType(txList []*blockchain.DefaultTransaction) []blockchain.Transaction {
	convTxList := make([]blockchain.Transaction, 0)

	for _, tx := range txList {
		convTxList = append(convTxList, tx)
	}

	return convTxList
}

func TestDefaultBlock_IsPrev(t *testing.T) {
	//given
	PrevBlock := blockchain.DefaultBlock{
		Seal: []byte("PrevSeal"),
	}

	b, err := PrevBlock.Serialize()
	assert.NoError(t, err)

	Block := blockchain.DefaultBlock{
		PrevSeal: []byte("PrevSeal"),
	}

	assert.Equal(t, true, Block.IsPrev(b))

}

func TestDefaultBlock_IsReadyToPublish(t *testing.T) {
	Block := blockchain.DefaultBlock{
		Seal: []byte("Seal"),
	}

	assert.Equal(t, true, Block.IsReadyToPublish())
}

func TestSerializeAndDeserialize(t *testing.T) {

	//given
	block := blockchain.DefaultBlock{
		Seal: []byte("Seal"),
	}

	//when
	serializedBlock, err := block.Serialize()

	//then
	assert.NoError(t, err)

	//given
	deserializedBlock := blockchain.DefaultBlock{}

	//when
	err = deserializedBlock.Deserialize(serializedBlock)

	//then
	assert.NoError(t, err)
	assert.Equal(t, deserializedBlock, block)
}

func TestSyncState(t *testing.T) {
	state := blockchain.NewSyncState(true, false)
	assert.Equal(t, state.SyncProgress(), true)
	assert.Equal(t, state.Sync(), false)

	state.SetSyncProgress(false)
	assert.Equal(t, state.SyncProgress(), false)

	state.SetSyncState(true)
	assert.Equal(t, state.Sync(), true)
}
