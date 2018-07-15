package api_gateway

import (
	"encoding/json"
	"errors"

	"log"

	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/leveldb-wrapper"
)

// this is an api only for querying current state which is repository of transaction
type TransactionQueryApi struct {
	transactionRepository TransactionPoolRepository
}

// find all transactions that are created by not committed as a block
func (t TransactionQueryApi) FindUncommittedTransactions() ([]txpool.Transaction, error) {

	return t.transactionRepository.FindAll()
}

// this repository is a current state of all uncommitted transactions
type TransactionPoolRepository interface {
	FindAll() ([]txpool.Transaction, error)
	Save(transaction txpool.Transaction) error
}

// this is an event_handler which listen all events related to transaction and update repository
// this struct will be relocated to other pkg
type TransactionEventListener struct {
	transactionRepository TransactionPoolRepository
}

// this function listens to TxCreatedEvent and update repository
func (t TransactionEventListener) HandleTransactionCreatedEvent(event txpool.TxCreatedEvent) {

	tx := event.GetTransaction()
	err := t.transactionRepository.Save(tx)

	if err != nil {
		log.Fatal(err.Error())
	}
}

type LeveldbTransactionPoolRepository struct {
	leveldb *leveldbwrapper.DB
}

func NewTransactionRepository(path string) *LeveldbTransactionPoolRepository {

	db := leveldbwrapper.CreateNewDB(path)
	db.Open()
	return &LeveldbTransactionPoolRepository{
		leveldb: db,
	}
}

func (t LeveldbTransactionPoolRepository) Save(transaction txpool.Transaction) error {

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

func (t LeveldbTransactionPoolRepository) Remove(id txpool.TransactionId) error {
	return t.leveldb.Delete([]byte(id), true)
}

func (t LeveldbTransactionPoolRepository) FindById(id txpool.TransactionId) (*txpool.Transaction, error) {
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

func (t LeveldbTransactionPoolRepository) FindAll() ([]txpool.Transaction, error) {

	iter := t.leveldb.GetIteratorWithPrefix([]byte(""))
	transactions := []txpool.Transaction{}

	for iter.Next() {
		val := iter.Value()
		tx := &txpool.Transaction{}
		err := txpool.Deserialize(val, tx)

		if err != nil {
			return nil, err
		}

		transactions = append(transactions, *tx)
	}

	return transactions, nil
}

func (t LeveldbTransactionPoolRepository) Close() {
	t.leveldb.Close()
}
