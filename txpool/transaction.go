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
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

const (
	VALID   TransactionStatus = 0
	INVALID TransactionStatus = 1

	General TransactionType = 0 + iota
)

type TransactionId = string

type TransactionStatus int
type TransactionType int

type TxData struct {
	ID        string
	Jsonrpc   string
	Method    string
	ICodeID   string
	Function  string
	Args      []string
	Signature []byte
}

//TxData Declaration
const (
	Invoke TxDataType = "invoke"
	Query  TxDataType = "query"
)

type TxDataType string

//Aggregate root must implement aggregate interface
type Transaction struct {
	ID        TransactionId
	Status    TransactionStatus
	TimeStamp time.Time
	Jsonrpc   string
	Method    TxDataType
	ICodeID   string
	Function  string
	Args      []string
	Signature []byte
	PeerID    string
}

// must implement id method
func (t Transaction) GetID() string {
	return string(t.ID)
}

// must implement on method
func (t *Transaction) On(transactionEvent midgard.Event) error {

	switch v := transactionEvent.(type) {

	case *event.TxCreated:

		t.ID = TransactionId(v.ID)
		t.PeerID = v.PeerID
		t.Status = TransactionStatus(v.Status)
		t.TimeStamp = v.TimeStamp
		t.Function = v.Function
		t.Args = v.Args
		t.Method = TxDataType(v.Method)
		t.Jsonrpc = v.Jsonrpc
		t.ICodeID = v.ICodeID
		t.Signature = t.Signature

	case *event.TxDeleted:
		t.ID = ""

	default:
		return errors.New(fmt.Sprintf("unhandled event [%s]", v))
	}

	return nil
}

func (t Transaction) Serialize() ([]byte, error) {
	return common.Serialize(t)
}

func Deserialize(b []byte, transaction *Transaction) error {

	err := json.Unmarshal(b, transaction)

	if err != nil {
		return err
	}

	return nil
}

func CreateTransaction(publisherId string, txData TxData) (Transaction, error) {

	id := xid.New().String()
	timeStamp := time.Now()

	transactionCreatedEvent := &event.TxCreated{
		EventModel: midgard.EventModel{
			ID:   id,
			Type: "transaction.created",
		},
		PeerID:    publisherId,
		Status:    int(VALID),
		TimeStamp: timeStamp,
		ICodeID:   txData.ICodeID,
		Jsonrpc:   txData.Jsonrpc,
		Method:    txData.Method,
		Signature: txData.Signature,
		Args:      txData.Args,
		Function:  txData.Function,
	}

	tx := &Transaction{}

	if err := saveAndOn(tx, transactionCreatedEvent); err != nil {
		return *tx, err
	}

	return *tx, nil
}

func DeleteTransaction(transaction Transaction) error {

	TxDeleteEvent := &event.TxDeleted{
		EventModel: midgard.EventModel{
			ID:   transaction.ID,
			Type: "transaction.deleted",
		},
	}

	return saveAndOn(&transaction, TxDeleteEvent)
}

//apply on aggrgate and publish to eventstore
func saveAndOn(aggregate midgard.Aggregate, event midgard.Event) error {

	//must do call on func first!!!
	//after save events if aggregate.On failed then data inconsistency will be occurred
	if err := aggregate.On(event); err != nil {
		return err
	}

	if err := eventstore.Save(event.GetID(), event); err != nil {
		return err
	}

	return nil
}
