package leveldb

import (
	"encoding/json"
	"errors"

	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/leveldb-wrapper"
)

type TransactionRepository struct {
	leveldb *leveldbwrapper.DB
}

func NewTransactionRepository(path string) *TransactionRepository {
	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return &TransactionRepository{
		leveldb: db,
	}
}

func (tr TransactionRepository) Save(transaction txpool.Transaction) error {
	if transaction.TxId.ToString() == "" {
		return errors.New("transaction ID is empty")
	}

	b, err := transaction.Serialize()

	if err != nil {
		return err
	}

	if err = tr.leveldb.Put([]byte(transaction.TxId), b, true); err != nil {
		return err
	}

	return nil
}

func (tr TransactionRepository) Remove(id txpool.TransactionId) error {
	return tr.leveldb.Delete([]byte(id), true)
}

func (tr TransactionRepository) FindById(id txpool.TransactionId) (*txpool.Transaction, error) {
	b, err := tr.leveldb.Get([]byte(id))

	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, nil
	}

	tx := &txpool.Transaction{}

	err = json.Unmarshal(b, tx)

	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (tr TransactionRepository) FindAll() ([]*txpool.Transaction, error) {

	iter := tr.leveldb.GetIteratorWithPrefix([]byte(""))
	transactions := []*txpool.Transaction{}
	for iter.Next() {
		val := iter.Value()
		tx := &txpool.Transaction{}
		err := txpool.Deserialize(val, tx)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}
