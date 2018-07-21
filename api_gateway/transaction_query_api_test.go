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
	"os"
	"testing"

	"time"

	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRepository_Save(t *testing.T) {
	// Given
	dbPath := "./.test"

	tr := NewTransactionRepository(dbPath)

	transaction := txpool.Transaction{
		ID:      "888",
		ICodeID: "889",
	}

	defer func() {
		tr.leveldb.Close()
		os.RemoveAll(dbPath)
	}()

	// When
	err := tr.Save(transaction)

	// Then
	t2 := &txpool.Transaction{}
	b, _ := tr.leveldb.Get([]byte(transaction.ID))
	txpool.Deserialize(b, t2)
	snapshot, _ := tr.leveldb.Snapshot()

	assert.Equal(t, transaction, *t2)
	assert.Equal(t, 1, len(snapshot))
	assert.Equal(t, nil, err)
}

func TestTransactionRepository_Remove(t *testing.T) {
	// Given
	dbPath := "./.test"

	tr := NewTransactionRepository(dbPath)

	transaction := txpool.Transaction{
		ID:      "888",
		ICodeID: "889",
	}

	defer func() {
		tr.leveldb.Close()
		os.RemoveAll(dbPath)
	}()

	// When
	_ = tr.Save(transaction)
	err := tr.Remove(transaction.ID)

	// Then
	t2 := &txpool.Transaction{}
	b, _ := tr.leveldb.Get([]byte(transaction.ID))
	txpool.Deserialize(b, t2)
	snapshot, _ := tr.leveldb.Snapshot()

	assert.Equal(t, txpool.Transaction{}, *t2)
	assert.Equal(t, 0, len(snapshot))
	assert.Equal(t, nil, err)

}

func TestTransactionRepository_FindById(t *testing.T) {
	// Given
	dbPath := "./.test"

	tr := NewTransactionRepository(dbPath)

	transaction := txpool.Transaction{
		ID:      "888",
		ICodeID: "889",
	}

	defer func() {
		tr.leveldb.Close()
		os.RemoveAll(dbPath)
	}()

	// When
	tr.Save(transaction)
	t2, err := tr.FindById("888")

	// Then
	assert.Equal(t, transaction, *t2)
	assert.Equal(t, nil, err)
}

func TestTransactionRepository_FindById_FindFailed(t *testing.T) {
	// Given
	dbPath := "./.test"

	tr := NewTransactionRepository(dbPath)

	defer func() {
		tr.leveldb.Close()
		os.RemoveAll(dbPath)
	}()

	// When
	t2, err := tr.FindById("888")

	// Then
	assert.Equal(t, true, t2 == nil)
	assert.Equal(t, nil, err)
}

func TestTransactionRepository_FindAll(t *testing.T) {
	// Given
	var transactions = []txpool.Transaction{}

	dbPath := "./.test"

	tr := NewTransactionRepository(dbPath)

	defer func() {
		tr.leveldb.Close()
		os.RemoveAll(dbPath)
	}()

	// When
	transactions1, err1 := tr.FindAll()

	// Then
	assert.Equal(t, transactions, transactions1)
	assert.Equal(t, nil, err1)

	// When
	t1 := txpool.Transaction{
		ID:      "888",
		ICodeID: "889",
	}
	tr.Save(t1)
	transactions2, err2 := tr.FindAll()

	// Then
	transactions = append(transactions, t1)

	assert.Equal(t, transactions, transactions2)
	assert.Equal(t, nil, err2)
}

func TestTransactionQueryApi_FindUncommittedTransactions(t *testing.T) {

	api, client, tearDown := setApiUp(t)

	defer tearDown()

	err := client.Publish("Event", "transaction.created", event.TxCreated{
		EventModel: midgard.EventModel{
			ID:      "1123",
			Time:    time.Now(),
			Type:    "transaction.created",
			Version: 3,
		},
		ICodeID: "2",
		Jsonrpc: "123",
	})

	assert.NoError(t, err)

	err = client.Publish("Event", "transaction.created", event.TxCreated{
		EventModel: midgard.EventModel{
			ID: "2",
		},
		ICodeID: "2",
		Jsonrpc: "123",
	})

	assert.NoError(t, err)

	err = client.Publish("Event", "transaction.created", event.TxCreated{
		EventModel: midgard.EventModel{
			ID: "3",
		},
		ICodeID: "2",
		Jsonrpc: "123",
	})

	assert.NoError(t, err)

	//wait until sync from event
	time.Sleep(3 * time.Second)

	txs, err := api.FindUncommittedTransactions()
	assert.Equal(t, len(txs), 3)
}

func TestTransactionEventListener_HandleTransactionDeletedEvent(t *testing.T) {

	api, client, tearDown := setApiUp(t)

	defer tearDown()

	err := client.Publish("Event", "transaction.created", event.TxCreated{
		EventModel: midgard.EventModel{
			ID:      "1123",
			Time:    time.Now(),
			Type:    "transaction.created",
			Version: 3,
		},
		ICodeID: "2",
		Jsonrpc: "123",
	})

	assert.NoError(t, err)

	//wait until sync from event
	time.Sleep(3 * time.Second)
	txs, err := api.FindUncommittedTransactions()
	assert.Equal(t, len(txs), 1)

	err = client.Publish("Event", "transaction.deleted", event.TxDeleted{
		EventModel: midgard.EventModel{
			ID: "1123",
		},
	})

	//wait until sync from event
	time.Sleep(3 * time.Second)

	txs, err = api.FindUncommittedTransactions()
	assert.Equal(t, len(txs), 0)
}

func setApiUp(t *testing.T) (TransactionQueryApi, *pubsub.Client, func()) {

	dbPath := "./.test"
	client := pubsub.Connect("")

	repo := NewTransactionRepository(dbPath)

	txQueryApi := TransactionQueryApi{transactionRepository: repo}
	txEventListener := &TransactionEventListener{transactionRepository: repo}

	err := client.Subscribe("Event", "transaction.*", txEventListener)
	assert.NoError(t, err)

	return txQueryApi, client, func() {
		os.RemoveAll(dbPath)
		client.Close()
	}
}
