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

package mem_test

import (
	"testing"

	"fmt"

	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/engine/txpool/infra/mem"
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

	for i := 0; i < 3; i++ {
		repo.Save(txpool.Transaction{
			ID: fmt.Sprintf("%s", i),
		})
	}

	//then
	transactionList, err := repo.FindAll()

	assert.NoError(t, err)

	i := 0
	for _, tx := range transactionList {
		assert.Equal(t, tx.ID, fmt.Sprintf("%s", i))
		i++
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
