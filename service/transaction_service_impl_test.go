package service

import (
	"testing"
	"os"
	"github.com/it-chain/it-chain-Engine/domain"
	"time"
	"github.com/stretchr/testify/assert"
	"github.com/it-chain/it-chain-Engine/common"
	"fmt"
)

func TestTransactionServiceImpl_AddTransaction(t *testing.T) {
	path := "./test"
	ts := NewTransactionService(path,nil,nil)
	defer func(){
		ts.Close()
		os.RemoveAll(path)
	}()

	tx := domain.CreateNewTransaction(
		"test",
		"test",
		0,
		time.Now().Round(0),
		&domain.TxData{})
	tx.GenerateHash()
	err := ts.AddTransaction(tx)
	assert.NoError(t, err)

	retrievedTX, err := ts.DB.GetDBHandle(WAITING_TRANSACTION).Get([]byte("test"))
	assert.NoError(t, err)

	deserializedTX := &domain.Transaction{}
	err = common.Deserialize(retrievedTX, deserializedTX)
	assert.NoError(t, err)
	assert.Equal(t, tx, deserializedTX)
}

func TestTransactionServiceImpl_DeleteTransactions(t *testing.T) {
	path := "./test"
	ts := NewTransactionService(path,nil,nil)
	defer func(){
		ts.Close()
		os.RemoveAll(path)
	}()

	tx := domain.CreateNewTransaction(
		"test",
		"test",
		0,
		time.Now().Round(0),
		&domain.TxData{})
	tx.GenerateHash()
	err := ts.AddTransaction(tx)
	assert.NoError(t, err)

	retrievedTX, err := ts.DB.GetDBHandle(WAITING_TRANSACTION).Get([]byte("test"))
	assert.NoError(t, err)

	deserializedTX := &domain.Transaction{}
	err = common.Deserialize(retrievedTX, deserializedTX)
	assert.NoError(t, err)
	assert.Equal(t, tx, deserializedTX)

	txs := make([]*domain.Transaction, 0)
	txs = append(txs, tx)
	ts.DeleteTransactions(txs)

	retrievedTX, err = ts.DB.GetDBHandle(WAITING_TRANSACTION).Get([]byte("test"))
	assert.NoError(t, err)
	assert.Equal(t, []byte(nil), retrievedTX)
}

func TestTransactionServiceImpl_GetTransactions(t *testing.T) {
	path := "./test"
	ts := NewTransactionService(path,nil,nil)
	defer func(){
		ts.Close()
		os.RemoveAll(path)
	}()

	txs := make([]*domain.Transaction, 0)

	for i := 1 ; i <= 2 ; i++ {
		tx := domain.CreateNewTransaction(
			"test",
			fmt.Sprintf("%d", i),
			0,
			time.Now().Round(0),
			&domain.TxData{})
		tx.GenerateHash()
		err := ts.AddTransaction(tx)
		assert.NoError(t, err)
		txs = append(txs, tx)
	}

	dbTXs, err := ts.GetTransactions(2)
	assert.NoError(t, err)
	assert.Equal(t, txs, dbTXs)
}