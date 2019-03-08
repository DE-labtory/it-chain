/*
 * Copyright 2018 DE-labtory
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

package mem

import (
	"errors"
	"sync"

	"github.com/DE-labtory/engine/txpool"
)

var ErrTransactionDoesNotExist = errors.New("transaction does not exist")
var ErrEmptyID = errors.New("transaction ID is empty")

func NewTransactionRepository() *TransactionRepository {
	return &TransactionRepository{
		TxMap:   make(map[txpool.TransactionId]txpool.Transaction),
		RWMutex: sync.RWMutex{},
	}
}

type TransactionRepository struct {
	TxMap map[txpool.TransactionId]txpool.Transaction
	sync.RWMutex
}

func (m *TransactionRepository) Save(transaction txpool.Transaction) error {

	m.Lock()
	defer m.Unlock()

	id := transaction.ID

	if id == "" {
		return ErrEmptyID
	}

	m.TxMap[id] = transaction

	return nil
}

func (m *TransactionRepository) Remove(id txpool.TransactionId) {
	m.Lock()
	defer m.Unlock()

	delete(m.TxMap, id)
}

func (m *TransactionRepository) FindById(id txpool.TransactionId) (txpool.Transaction, error) {
	m.Lock()
	defer m.Unlock()

	t, ok := m.TxMap[id]

	if ok {
		return t, nil
	}

	return txpool.Transaction{}, ErrTransactionDoesNotExist
}

func (m *TransactionRepository) FindAll() ([]txpool.Transaction, error) {
	m.Lock()
	defer m.Unlock()

	s := make([]txpool.Transaction, 0)

	for _, transaction := range m.TxMap {
		s = append(s, transaction)
	}

	return s, nil
}
