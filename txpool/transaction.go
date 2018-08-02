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

package txpool

import (
	"time"

	"github.com/rs/xid"
)

type TransactionId = string

type TxData struct {
	Jsonrpc   string
	ICodeID   string
	Function  string
	Args      []string
	Signature []byte
}

//Aggregate root must implement aggregate interface
type Transaction struct {
	ID        TransactionId
	TimeStamp time.Time
	Jsonrpc   string
	ICodeID   string
	Function  string
	Args      []string
	Signature []byte
	PeerID    string
}

func CreateTransaction(publisherId string, txData TxData) (Transaction, error) {

	id := xid.New().String()
	timeStamp := time.Now()

	transaction := Transaction{
		ID:        id,
		PeerID:    publisherId,
		TimeStamp: timeStamp,
		ICodeID:   txData.ICodeID,
		Jsonrpc:   txData.Jsonrpc,
		Signature: txData.Signature,
		Args:      txData.Args,
		Function:  txData.Function,
	}

	return transaction, nil
}

type TransactionRepository interface {
	FindAll() ([]Transaction, error)
	Save(transaction Transaction) error
	Remove(id TransactionId)
	FindById(id TransactionId) (Transaction, error)
}
