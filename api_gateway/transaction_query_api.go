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
	"encoding/json"
	"errors"
	"log"

	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/leveldb-wrapper"
)

// this is an api only for querying current state which is repository of transaction
type TransactionQueryApi struct {
	transactionRepository TransactionPoolRepository
}

func NewTransactionQueryApi(transactionRepository TransactionPoolRepository) TransactionQueryApi {
	return TransactionQueryApi{
		transactionRepository: transactionRepository,
	}
}

// find all transactions that are created by not committed as a block
func (t TransactionQueryApi) FindUncommittedTransactions() ([]txpool.Transaction, error) {

	return t.transactionRepository.FindAll()
}

// this repository is a current state of all uncommitted transactions
type TransactionPoolRepository interface {
	FindAll() ([]txpool.Transaction, error)
	Save(transaction txpool.Transaction) error
	Remove(id txpool.TransactionId) error
}

// this is an event_handler which listen all events related to transaction and update repository
// this struct will be relocated to other pkg
type TransactionEventListener struct {
	transactionRepository TransactionPoolRepository
}

func NewTransactionEventListener(transactionRepository TransactionPoolRepository) TransactionEventListener {
	return TransactionEventListener{
		transactionRepository: transactionRepository,
	}
}

// this function listens to TxCreatedEvent and update repository
func (t TransactionEventListener) HandleTransactionCreatedEvent(event txpool.TxCreatedEvent) {

	tx := event.GetTransaction()
	err := t.transactionRepository.Save(tx)

	if err != nil {
		log.Fatal(err.Error())
	}
}

func (t TransactionEventListener) HandleTransactionDeletedEvent(event txpool.TxDeletedEvent) {

	err := t.transactionRepository.Remove(event.GetID())

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
