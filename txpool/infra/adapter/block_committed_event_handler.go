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

package adapter

import (
	"errors"

	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/engine/txpool/api"
)

var ErrNoEventID = errors.New("no event id ")

type BlockCommittedEventHandler struct {
	transactionApi api.TransactionApi
}

func (e BlockCommittedEventHandler) HandleBlockCommittedEvent(event txpool.BlockCommittedEvent) error {

	txs := event.Transactions

	for _, tx := range txs {
		err := e.transactionApi.DeleteTransaction(tx.TxId)

		if err != nil {
			return err
		}
	}

	return nil
}
