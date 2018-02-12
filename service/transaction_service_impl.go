package service

import (
	"it-chain/db/leveldbhelper"
	"it-chain/domain"
	"it-chain/common"
)

const (
	WAITING_TRANSACTION = "waiting_transaction"
)

type TransactionServiceImpl struct {
	DB *leveldbhelper.DBProvider
}

func CreateNewTransactionService(path string) *TransactionServiceImpl {
	return &TransactionServiceImpl{DB: leveldbhelper.CreateNewDBProvider(path)}
}

func (t *TransactionServiceImpl) Close() {
	t.DB.Close()
}

func (t *TransactionServiceImpl) AddTransaction(tx *domain.Transaction) error {
	db := t.DB.GetDBHandle(WAITING_TRANSACTION)
	serializedTX, err := common.Serialize(tx)
	if err != nil {
		return err
	}

	err = db.Put([]byte(tx.TransactionID), serializedTX, true)
	return err
}

func (t *TransactionServiceImpl) DeleteTransactions(txs []*domain.Transaction) error {
	db := t.DB.GetDBHandle(WAITING_TRANSACTION)
	batch := make(map[string][]byte)

	for _, tx := range txs {
		batch[tx.TransactionID] = nil
	}

	return db.WriteBatch(batch, true)
}

func (t *TransactionServiceImpl) GetTransactions() ([]*domain.Transaction, error) {
	db := t.DB.GetDBHandle(WAITING_TRANSACTION)
	iter := db.GetIteratorWithPrefix()
	ret := make([]*domain.Transaction, 0)

	for iter.Next() {
		val := iter.Value()
		tx := &domain.Transaction{}
		err := common.Deserialize(val, tx)

		if err != nil {
			return nil, err
		}

		ret = append(ret, tx)
	}

	iter.Release()

	return ret, nil
}