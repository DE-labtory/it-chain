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
	"github.com/it-chain/engine/common/command"
)

type BlockProposeApi interface {
	CreateProposedBlock(txList []*blockchain.DefaultTransaction) (blockchain.DefaultBlock, error)
	ConsentBlock(consensusType string, block blockchain.DefaultBlock) error
}

type BlockProposeCommandHandler struct {
	blockApi   BlockProposeApi
	engineMode string
}

func NewBlockProposeCommandHandler(blockApi BlockProposeApi, engineMode string) *BlockProposeCommandHandler {
	return &BlockProposeCommandHandler{
		blockApi:   blockApi,
		engineMode: engineMode,
	}
}

func (h *BlockProposeCommandHandler) HandleProposeBlockCommand(command command.ProposeBlock) error {

	if err := validateCommand(command); err != nil {
		return err
	}

	txList := command.TxList
	defaultTxList := getBackTxList(txList)

	proposedBlock, err := h.blockApi.CreateProposedBlock(defaultTxList)
	if err != nil {
		return err
	}

	if err := h.blockApi.ConsentBlock(h.engineMode, proposedBlock); err != nil {
		return err
	}

	return nil
}

func validateCommand(command command.ProposeBlock) error {
	txList := command.TxList

	if txList == nil || len(txList) == 0 {
		return ErrCommandTransactions
	}
	return nil
}

func getBackTxList(txList []command.Tx) []*blockchain.DefaultTransaction {
	defaultTxList := make([]*blockchain.DefaultTransaction, 0)

	for _, tx := range txList {
		defaultTx := getBackTx(tx)
		defaultTxList = append(defaultTxList, defaultTx)
	}
	return defaultTxList
}

func getBackTx(tx command.Tx) *blockchain.DefaultTransaction {
	return &blockchain.DefaultTransaction{
		ID:        tx.ID,
		ICodeID:   tx.ICodeID,
		PeerID:    tx.PeerID,
		Timestamp: tx.TimeStamp,
		Jsonrpc:   tx.Jsonrpc,
		Function:  tx.Function,
		Args:      tx.Args,
		Signature: tx.Signature,
	}
}
