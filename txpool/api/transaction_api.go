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

package api

import (
	"log"

	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/engine/txpool"
)

type TransactionApi struct {
	publisherId string
}

func NewTransactionApi(publisherId string) TransactionApi {
	return TransactionApi{
		publisherId: publisherId,
	}
}

func (t TransactionApi) CreateTransaction(txData txpool.TxData) (txpool.Transaction, error) {

	log.Printf("create transaction: [%v]", txData)

	tx, err := txpool.CreateTransaction(t.publisherId, txData)

	if err != nil {
		log.Printf("fail to transaction: [%v]", err)
		return tx, err
	}

	log.Printf("transaction created: [%v]", tx)
	return tx, nil
}

func (t TransactionApi) DeleteTransaction(id txpool.TransactionId) error {

	log.Printf("delete transaction: [%v]", id)

	tx := &txpool.Transaction{}

	if err := eventstore.Load(tx, id); err != nil {
		log.Printf("fail to delete transaction: [%v]", id)
		return err
	}

	return txpool.DeleteTransaction(*tx)
}
