package api_gateway

import (
	"os"
	"testing"

	"time"

	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/midgard"
	"github.com/it-chain/midgard/bus/rabbitmq"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRepository_Save(t *testing.T) {
	// Given
	dbPath := "./.test"

	tr := NewTransactionRepository(dbPath)

	transaction := txpool.Transaction{
		TxId: "888",
		TxData: txpool.TxData{
			ID: "889",
		},
	}

	defer func() {
		tr.leveldb.Close()
		os.RemoveAll(dbPath)
	}()

	// When
	err := tr.Save(transaction)

	// Then
	t2 := &txpool.Transaction{}
	b, _ := tr.leveldb.Get([]byte(transaction.TxId))
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
		TxId: "888",
		TxData: txpool.TxData{
			ID: "889",
		},
	}

	defer func() {
		tr.leveldb.Close()
		os.RemoveAll(dbPath)
	}()

	// When
	_ = tr.Save(transaction)
	err := tr.Remove(transaction.TxId)

	// Then
	t2 := &txpool.Transaction{}
	b, _ := tr.leveldb.Get([]byte(transaction.TxId))
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
		TxId: "888",
		TxData: txpool.TxData{
			ID: "889",
		},
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
	var transactions = []*txpool.Transaction{}

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
		TxId: "888",
		TxData: txpool.TxData{
			ID: "889",
		},
	}
	tr.Save(t1)
	transactions2, err2 := tr.FindAll()

	// Then
	transactions = append(transactions, &t1)

	assert.Equal(t, transactions, transactions2)
	assert.Equal(t, nil, err2)
}

func TestTransactionQueryApi_FindUncommittedTransactions(t *testing.T) {

	api, client, tearDown := setApiUp(t)

	defer tearDown()

	err := client.Publish("Event", "transaction.created", txpool.TxCreatedEvent{
		EventModel: midgard.EventModel{
			ID:      "1123",
			Time:    time.Now(),
			Type:    "transaction.created",
			Version: 3,
		},
		ID:      "1",
		ICodeID: "2",
		TxHash:  "123",
		Jsonrpc: "123",
	})

	assert.NoError(t, err)

	err = client.Publish("Event", "transaction.created", txpool.TxCreatedEvent{
		EventModel: midgard.EventModel{
			ID: "2",
		},
		ID:      "2",
		ICodeID: "2",
		TxHash:  "123",
		Jsonrpc: "123",
	})

	assert.NoError(t, err)

	err = client.Publish("Event", "transaction.created", txpool.TxCreatedEvent{
		EventModel: midgard.EventModel{
			ID: "3",
		},
		ID:      "3",
		ICodeID: "2",
		TxHash:  "123",
		Jsonrpc: "123",
	})

	assert.NoError(t, err)

	txs, err := api.FindUncommittedTransactions()

	assert.Equal(t, len(txs), 3)
}

func setApiUp(t *testing.T) (TransactionQueryApi, *rabbitmq.Client, func()) {

	dbPath := "./.test"
	client := rabbitmq.Connect("")

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
