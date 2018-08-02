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

package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"time"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/event"
	ygg "github.com/it-chain/yggdrasill/common"
)

// Status 변수는 Transaction의 상태를 Unconfirmed, Confirmed, Unknown 중 하나로 표현함.
type Status int

type Transaction = ygg.Transaction

// DefaultTransaction 구조체는 Transaction 인터페이스의 기본 구현체이다.
type DefaultTransaction struct {
	ID        string
	Function  string
	Args      []string
	PeerID    string
	Jsonrpc   string
	ICodeID   string
	Timestamp time.Time
	Signature []byte
}

// GetID 함수는 Transaction의 ID 값을 반환한다.
func (t *DefaultTransaction) GetID() string {
	return t.ID
}

func (t *DefaultTransaction) GetContent() ([]byte, error) {
	content := struct {
		ID        string
		PeerID    string
		Timestamp time.Time
	}{t.ID,t.PeerID, t.Timestamp}

	serialized, err := serialize(content)
	if err != nil {
		return nil, err
	}

	return serialized, nil
}

func (t *DefaultTransaction) GetSignature() []byte {
	return t.Signature
}

// CalculateSeal 함수는 Transaction 고유의 Hash 값을 계산하여 반환한다.
func (t *DefaultTransaction) CalculateSeal() ([]byte, error) {
	serializedTx, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	return calculateHash(serializedTx), nil
}

func (t *DefaultTransaction) SetSignature(signature []byte) {
	t.Signature = signature
}

// Serialize 함수는 Transaction을 []byte 형태로 변환한다.
func (t *DefaultTransaction) Serialize() ([]byte, error) {
	return serialize(t)
}

func (t *DefaultTransaction) Deserialize(serializedBytes []byte) error {
	if len(serializedBytes) == 0 {
		return nil
	}

	err := json.Unmarshal(serializedBytes, t)

	if err != nil {
		return err
	}

	return nil
}

// NewDefaultTransaction 함수는 새로운 DefaultTransaction를 반환한다.
func NewDefaultTransaction(peerID string, txID string, timestamp time.Time) *DefaultTransaction {
	return &DefaultTransaction{
		ID:        txID,
		PeerID:    peerID,
		Timestamp: timestamp,
	}
}

func serialize(data interface{}) ([]byte, error) {
	serialized, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	return serialized, nil
}

func calculateHash(b []byte) []byte {
	hashValue := sha256.New()
	hashValue.Write(b)
	return hashValue.Sum(nil)
}

func deserializeTxList(txList []byte) ([]*DefaultTransaction, error) {
	DefaultTxList := []*DefaultTransaction{}

	err := common.Deserialize(txList, &DefaultTxList)

	if err != nil {
		return nil, err
	}
	return DefaultTxList, nil
}

func ConvertToTransactionList(txList []event.Tx) []*DefaultTransaction {
	defaultTxList := make([]*DefaultTransaction, 0)

	for _, tx := range txList {
		defaultTx := ConvertToTransaction(tx)
		defaultTxList = append(defaultTxList, defaultTx)
	}

	return defaultTxList
}

func ConvertToTransaction(tx event.Tx) *DefaultTransaction {
	return &DefaultTransaction{
		ID:        tx.ID,
		PeerID:    tx.PeerID,
		ICodeID:   tx.ICodeID,
		Timestamp: tx.TimeStamp,
		Signature: tx.Signature,
	}
}

func ConvBackFromTransactionList(defaultTxList []*DefaultTransaction) []event.Tx {
	txList := make([]event.Tx, 0)

	for _, defaultTx := range defaultTxList {
		tx := ConvBackFromTransaction(defaultTx)
		txList = append(txList, tx)
	}

	return txList
}

func ConvBackFromTransaction(defaultTx *DefaultTransaction) event.Tx {
	return event.Tx{
		ID:        defaultTx.ID,
		ICodeID:   defaultTx.ICodeID,
		TimeStamp: defaultTx.Timestamp,
		Jsonrpc:   defaultTx.Jsonrpc,
		Function:  defaultTx.Function,
		Args:      defaultTx.Args,
		Signature: defaultTx.Signature,
	}
}

func ConvertTxType(txList []*DefaultTransaction) []Transaction {
	convTxList := make([]Transaction, 0)

	for _, tx := range txList {
		convTxList = append(convTxList, tx)
	}

	return convTxList
}

func GetBackTxType(txList []Transaction) []*DefaultTransaction {
	convTxList := make([]*DefaultTransaction, 0)

	for _, tx := range txList {
		convTxList = append(convTxList, tx.(*DefaultTransaction))
	}

	return convTxList
}
