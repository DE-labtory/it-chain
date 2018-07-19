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

package memory

import (
	"sync"

	"errors"

	"github.com/it-chain/engine/txpool"
)

var ErrAlreadyExist = errors.New("transaction already exist")
var ErrDoesNotExist = errors.New("transaction does not exist")

type MemoryTransactionRepository struct {
	transactions map[txpool.TransactionId]txpool.Transaction
	mux          sync.RWMutex
}

func NewMemoryTransactionRepository() MemoryTransactionRepository {
	return MemoryTransactionRepository{
		transactions: make(map[txpool.TransactionId]txpool.Transaction),
		mux:          sync.RWMutex{},
	}
}

func (m MemoryTransactionRepository) Save(transaction txpool.Transaction) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	m.transactions[transaction.TxId] = transaction

	return nil
}

func (m MemoryTransactionRepository) Remove(id txpool.TransactionId) error {
	m.mux.Lock()
	defer m.mux.Unlock()

	delete(m.transactions, id)

	return nil
}

func (m MemoryTransactionRepository) FindById(id txpool.TransactionId) (txpool.Transaction, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	tx, ok := m.transactions[id]

	if !ok {
		return txpool.Transaction{}, ErrDoesNotExist
	}

	return tx, nil
}

func (m MemoryTransactionRepository) FindAll() ([]txpool.Transaction, error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	v := make([]txpool.Transaction, 0, len(m.transactions))

	for _, value := range m.transactions {
		v = append(v, value)
	}

	return v, nil
}
