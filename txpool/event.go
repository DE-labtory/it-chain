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

	"github.com/it-chain/midgard"
)

type TxCreatedEvent struct {
	midgard.EventModel
	PublishPeerId string
	TxStatus      int
	TxHash        string
	TimeStamp     time.Time
	Jsonrpc       string
	Method        string
	Params        Param
	ICodeID       string
}

func (tx TxCreatedEvent) GetTransaction() Transaction {

	return Transaction{
		TxId:          TransactionId(tx.EventModel.ID),
		PublishPeerId: tx.PublishPeerId,
		TxStatus:      TransactionStatus(tx.TxStatus),
		TxHash:        tx.TxHash,
		TxData: TxData{
			ICodeID: tx.ICodeID,
			Jsonrpc: tx.Jsonrpc,
			Method:  TxDataType(tx.Method),
			Params:  tx.Params,
		},
		TimeStamp: tx.TimeStamp,
	}
}

// when block committed check transaction and delete
type TxDeletedEvent struct {
	midgard.EventModel
}

type BlockCommittedEvent struct {
	midgard.EventModel
	Transactions []Transaction
}
