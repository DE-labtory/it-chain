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
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/txpool"
	"errors"
	"log"
)

var ErrCommandTransactions = errors.New("command's transactions nil or have length of zero")
var ErrTxHasMissingProperties = errors.New("Tx has missing properties")

type BlockCreateApi interface {
	CreateBlock(txList []blockchain.Transaction) error
}

type BlockProposeCommandHandler struct{
	blockApi BlockCreateApi
}

func NewBlockProposeCommandHandler(blockApi BlockCreateApi) *BlockProposeCommandHandler{
	return &BlockProposeCommandHandler{
		blockApi: blockApi,
	}
}

func (h *BlockProposeCommandHandler) HandleProposeBlockCommand(command blockchain.ProposeBlockCommand) {
	if err := validateCommand(command); err != nil {
		log.Fatal(err)
		return
	}
	txList := command.Transactions

	if err := validateTxList(txList); err != nil {
		log.Fatal(err)
		return
	}

	defaultTxList := convertTxList(txList)

	if err := h.blockApi.CreateBlock(defaultTxList); err != nil {
		log.Fatal(err)
	}
}

func validateCommand(command blockchain.ProposeBlockCommand) error {
	txList := command.Transactions

	if txList == nil || len(txList) == 0 {
		return ErrCommandTransactions
	}
	return nil
}

func validateTxList(txList []txpool.Transaction) error {
	var err error

	for _, tx := range txList {
		err = validateTx(tx)
	}

	return err
}

func validateTx(tx txpool.Transaction) error {
	if tx.TxId == "" || tx.PublishPeerId == "" || tx.TimeStamp.IsZero() || tx.TxData.Jsonrpc == "" ||
		tx.TxData.Method == "" || tx.TxData.Params.Function == "" || tx.TxData.Params.Args == nil {
		return ErrTxHasMissingProperties
	}
	return nil
}

func convertTxList(txList []txpool.Transaction) []blockchain.Transaction {
	defaultTxList := make([]blockchain.Transaction, 0)

	for _, tx := range txList {
		defaultTx := convertTx(tx)
		defaultTxList = append(defaultTxList, defaultTx)
	}

	return defaultTxList
}

func convertTx(tx txpool.Transaction) blockchain.Transaction {
	return &blockchain.DefaultTransaction{
		ID: tx.GetID(),
		Status: blockchain.Status(tx.TxStatus),
		PeerID: tx.PublishPeerId,
		Timestamp: tx.TimeStamp,
		TxData: &blockchain.TxData{
			Jsonrpc: tx.TxData.Jsonrpc,
			Method: blockchain.TxDataType(tx.TxData.Method),
			Params: blockchain.Params{
				Function: tx.TxData.Params.Function,
				Args: tx.TxData.Params.Args,
			},
		},
	}
}