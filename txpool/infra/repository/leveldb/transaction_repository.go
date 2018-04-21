package leveldb

import (
	"encoding/json"
	"errors"

	"github.com/it-chain/it-chain-Engine/txpool/domain/model/transaction"
	"github.com/it-chain/leveldb-wrapper"
)

//todo generate test code
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

func (tr TransactionRepository) Save(transaction transaction.Transaction) error {
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

func (tr TransactionRepository) Remove(id transaction.TransactionId) error {
	return tr.leveldb.Delete([]byte(id), true)
}

func (tr TransactionRepository) FindById(id transaction.TransactionId) (*transaction.Transaction, error) {
	b, err := tr.leveldb.Get([]byte(id))

	if err != nil {
		return nil, err
	}

	if len(b) == 0 {
		return nil, nil
	}

	tx := &transaction.Transaction{}

	err = json.Unmarshal(b, tx)

	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (tr TransactionRepository) FindAll() ([]*transaction.Transaction, error) {
	iter := tr.leveldb.GetIteratorWithPrefix([]byte(""))
	transactions := []*transaction.Transaction{}
	for iter.Next() {
		val := iter.Value()
		tx := &transaction.Transaction{}
		err := transaction.Deserialize(val, tx)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}
