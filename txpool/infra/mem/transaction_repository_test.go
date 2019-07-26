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

package mem_test

import (
	"strconv"
	"testing"

	"github.com/DE-labtory/it-chain/txpool"
	"github.com/DE-labtory/it-chain/txpool/infra/mem"
	"github.com/stretchr/testify/assert"
)

func TestTransactionRepository_Save(t *testing.T) {

	tests := map[string]struct {
		input txpool.Transaction
		err   error
	}{
		"success": {
			input: txpool.Transaction{
				ID:       "1",
				Function: "initA",
				Jsonrpc:  "2.0",
			},
			err: nil,
		},
		"fail empty id": {
			input: txpool.Transaction{
				Function: "initA",
				Jsonrpc:  "2.0",
			},
			err: mem.ErrEmptyID,
		},
	}

	repo := mem.NewTransactionRepository()

	for _, test := range tests {
		err := repo.Save(test.input)
		assert.Equal(t, err, test.err)
	}
}

func TestTransactionRepository_FindAll(t *testing.T) {

	//given
	repo := mem.NewTransactionRepository()

	testTransactions := make([]txpool.Transaction, 0)

	for i := 0; i < 3; i++ {
		testTransactions = append(testTransactions, txpool.Transaction{
			ID: strconv.Itoa(i),
		})
	}

	for _, tx := range testTransactions {
		repo.Save(tx)
	}

	//then
	transactionList, err := repo.FindAll()

	assert.NoError(t, err)

	for _, tx := range testTransactions {
		assert.Contains(t, transactionList, tx)
	}
}

func TestTransactionRepository_Remove(t *testing.T) {

	//given
	repo := mem.NewTransactionRepository()
	repo.Save(txpool.Transaction{
		ID: "123",
	})

	//when
	repo.Remove("123")

	//then
	_, err := repo.FindById("123")
	assert.Equal(t, err, mem.ErrTransactionDoesNotExist)
}
