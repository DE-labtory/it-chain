package leveldb

import (
	"os"
	"testing"

	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/magiconair/properties/assert"
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
