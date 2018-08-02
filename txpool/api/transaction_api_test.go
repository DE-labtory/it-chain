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

package api_test

import (
	"testing"

	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/engine/txpool/api"
	"github.com/it-chain/engine/txpool/infra/mem"
	"github.com/stretchr/testify/assert"
)

func TestTransactionApi_CreateTransaction(t *testing.T) {

	tests := map[string]struct {
		input struct {
			txData txpool.TxData
		}
		err error
	}{
		"success": {
			input: struct {
				txData txpool.TxData
			}{
				txData: txpool.TxData{
					ICodeID:   "gg",
					Function:  "1",
					Signature: []byte("123"),
					Args:      []string{"1", "2"},
					Jsonrpc:   "2.0",
				},
			},
			err: nil,
		},
	}

	transactionRepository := mem.NewTransactionRepository()
	transactionApi := api.NewTransactionApi("zf", transactionRepository)

	for _, test := range tests {
		tx, err := transactionApi.CreateTransaction(test.input.txData)

		assert.Equal(t, test.err, err)
		assert.Equal(t, tx.ICodeID, test.input.txData.ICodeID)
		assert.Equal(t, tx.Args, test.input.txData.Args)
		assert.Equal(t, tx.Signature, test.input.txData.Signature)
		assert.Equal(t, tx.Jsonrpc, test.input.txData.Jsonrpc)
		assert.Equal(t, tx.Function, test.input.txData.Function)
	}
}

func TestTransactionApi_DeleteTransaction(t *testing.T) {

	tests := map[string]struct {
		input string
		err   error
	}{
		"success": {
			input: "transactionID",
			err:   mem.ErrTransactionDoesNotExist,
		},
	}

	transactionRepository := mem.NewTransactionRepository()
	transactionApi := api.NewTransactionApi("zf", transactionRepository)

	transactionRepository.Save(txpool.Transaction{
		ID: "transactionID",
	})

	for _, test := range tests {
		transactionApi.DeleteTransaction(test.input)

		_, err := transactionRepository.FindById(test.input)
		assert.Equal(t, err, test.err)
	}
}
