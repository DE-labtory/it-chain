package api_gateway

import (
	"encoding/json"
	"errors"

	"github.com/it-chain/it-chain-Engine/txpool"
	leveldbwrapper "github.com/it-chain/leveldb-wrapper"
)

// this is an api only for querying current state which is repository of transaction
type TransactionQueryApi struct {
	transactionRepository TransactionPoolRepository
}

// find all transactions that are created by not committed as a block
func (t TransactionQueryApi) FindUncommittedTransactions() []txpool.Transaction {
	return t.transactionRepository.FindAllTransaction()
}

// this repository is a current state of all uncommitted transactions
type TransactionPoolRepository interface {
	FindAllTransaction() []txpool.Transaction
	Save(transaction txpool.Transaction)
}

// this is an event_handler which listen all events related to transaction and update repository
// this struct will be relocated to other pkg
type TransactionEventListener struct {
	transactionRepository TransactionPoolRepository
}

// this function listens to TxCreatedEvent and update repository
func (t TransactionEventListener) HandleTransactionCreatedEvent(event txpool.TxCreatedEvent) {

	tx := event.GetTransaction()
	t.transactionRepository.Save(tx)
}

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

func (t TransactionRepository) Save(transaction txpool.Transaction) error {
	if transaction.TxId == "" {
		return errors.New("transaction ID is empty")
	}

	b, err := transaction.Serialize()

	if err != nil {
		return err
	}

	if err = t.leveldb.Put([]byte(transaction.TxId), b, true); err != nil {
		return err
	}

	return nil
}

func (t TransactionRepository) Remove(id txpool.TransactionId) error {
	return t.leveldb.Delete([]byte(id), true)
}

func (t TransactionRepository) FindById(id txpool.TransactionId) (*txpool.Transaction, error) {
	b, err := t.leveldb.Get([]byte(id))

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

func (t TransactionRepository) FindAll() ([]*txpool.Transaction, error) {

	iter := t.leveldb.GetIteratorWithPrefix([]byte(""))
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

func (t TransactionRepository) Close() {
	t.leveldb.Close()
}
