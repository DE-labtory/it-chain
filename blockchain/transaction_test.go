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

package blockchain_test

import (
	"testing"
	"time"

	"github.com/DE-labtory/engine/blockchain"
	"github.com/DE-labtory/engine/common/event"
	"github.com/stretchr/testify/assert"
)

func TestDefaultTransaction_SetGetFunctions(t *testing.T) {
	tx := blockchain.DefaultTransaction{
		ID:        "tx01",
		ICodeID:   "ICode01",
		PeerID:    "Peer01",
		Timestamp: time.Now().Round(0),
		Jsonrpc:   "json01",
		Function:  "function01",
		Args:      nil,
		Signature: nil,
	}

	//when
	txID := tx.GetID()

	//then
	assert.Equal(t, "tx01", txID)

	//when
	content, err := tx.GetContent()
	assert.NoError(t, err)

	//then
	b, err := tx.Serialize()
	assert.NoError(t, err)
	assert.Equal(t, b, content)

	//when
	tx.SetSignature([]byte("Signature"))

	//then
	assert.Equal(t, []byte("Signature"), tx.GetSignature())

}

func TestDefaultTransaction_CalculateSeal(t *testing.T) {
	TimeStamp := getTestingTime()

	tx := blockchain.DefaultTransaction{
		ID:        "tx01",
		ICodeID:   "ICode01",
		PeerID:    "Peer01",
		Timestamp: TimeStamp,
		Jsonrpc:   "json01",
		Function:  "function01",
		Args:      nil,
		Signature: []byte("signature"),
	}

	b := []byte{15, 117, 178, 0, 211, 233, 104, 240, 164, 241, 236, 75, 95, 51, 237, 243, 5, 179, 182, 13, 17, 76, 247, 50, 150, 212, 73, 64, 14, 3, 71, 4}

	seal, err := tx.CalculateSeal()
	assert.NoError(t, err)
	assert.Equal(t, b, seal)
}

func TestDefaultTransaction_SerializeAndDeserialize(t *testing.T) {
	//given
	tx := blockchain.DefaultTransaction{
		ID:        "tx01",
		ICodeID:   "ICode01",
		PeerID:    "Peer01",
		Timestamp: time.Now().Round(0),
		Jsonrpc:   "json01",
		Function:  "function01",
		Args:      nil,
		Signature: []byte("signature"),
	}

	//when
	txBytes, err := tx.Serialize()

	//then
	assert.NoError(t, err)

	//given
	deserializedTx := blockchain.DefaultTransaction{}

	//when
	err = deserializedTx.Deserialize(txBytes)

	//then
	assert.NoError(t, err)
	assert.Equal(t, deserializedTx.ID, tx.ID)
}

func TestConvertTxList(t *testing.T) {
	//given
	EventTxList1 := []event.Tx{
		{
			ID:        "tx01",
			ICodeID:   "ICode01",
			PeerID:    "Peer01",
			TimeStamp: time.Now().Round(0),
			Jsonrpc:   "json01",
			Function:  "function01",
			Args:      nil,
			Signature: []byte("signature"),
		},
	}

	//when
	DefaultTxList1 := blockchain.ConvertToTransactionList(EventTxList1)

	//then
	assert.Equal(t, "tx01", DefaultTxList1[0].GetID())

	//when
	EventTxList2 := blockchain.ConvBackFromTransactionList(DefaultTxList1)

	//then
	assert.Equal(t, EventTxList1, EventTxList2)

	//when
	TxList := blockchain.ConvertTxType(DefaultTxList1)
	DefaultTxList2 := blockchain.GetBackTxType(TxList)

	assert.Equal(t, DefaultTxList1, DefaultTxList2)
}

func getTestingTime() time.Time {
	const longForm = "Jan 2, 2006 at 3:04pm (MST)"
	testingTime, _ := time.Parse(longForm, "Feb 3, 2013 at 7:54pm (UTC)")

	return testingTime
}
