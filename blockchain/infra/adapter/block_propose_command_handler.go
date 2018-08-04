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
	"github.com/it-chain/engine/common/rabbitmq/rpc"
)

type BlockCreateApi interface {
	CreateProposedBlock(txList []blockchain.Transaction) (blockchain.DefaultBlock, error)
}

type BlockSaveRepository interface {
	Save(block blockchain.DefaultBlock) error
}

type BlockCommitService interface {
	CommitBlock(block blockchain.DefaultBlock) error
}

type BlockProposeCommandHandler struct {
	blockApi        BlockCreateApi
	blockRepository BlockSaveRepository
	eventService    BlockCommitService
	engineMode      string
}

func NewBlockProposeCommandHandler(blockApi BlockCreateApi, blockRepository BlockSaveRepository, eventService BlockCommitService, engineMode string) *BlockProposeCommandHandler {
	return &BlockProposeCommandHandler{
		blockApi:        blockApi,
		blockRepository: blockRepository,
		eventService:    eventService,
		engineMode:      engineMode,
	}
}

func (h *BlockProposeCommandHandler) HandleProposeBlockCommand(command command.ProposeBlock) (struct{}, rpc.Error) {

	if err := validateCommand(command); err != nil {
		return struct{}{}, rpc.Error{Message: err.Error()}
	}

	txList := command.TxList

	if err := validateTxList(txList); err != nil {
		return struct{}{}, rpc.Error{Message: err.Error()}
	}

	defaultTxList := convertTxList(txList)

	if h.engineMode == "solo" {

		//create
		ProposedBlock, err := h.blockApi.CreateProposedBlock(defaultTxList)

		if err != nil {
			return struct{}{}, rpc.Error{Message: err.Error()}
		}

		// save
		err = h.blockRepository.Save(ProposedBlock)
		if err != nil {
			return struct{}{}, rpc.Error{Message: err.Error()}
		}

		// event
		err = h.eventService.CommitBlock(ProposedBlock)

		if err != nil {
			return struct{}{}, rpc.Error{Message: err.Error()}
		}
		return struct{}{}, rpc.Error{}
	}

	return struct{}{}, rpc.Error{}
}

func validateCommand(command command.ProposeBlock) error {
	txList := command.TxList

	if txList == nil || len(txList) == 0 {
		return ErrCommandTransactions
	}
	return nil
}

func validateTxList(txList []command.Tx) error {
	var err error

	for _, tx := range txList {
		err = validateTx(tx)
	}

	return err
}

func validateTx(tx command.Tx) error {
	if tx.ID == "" || tx.PeerID == "" || tx.TimeStamp.IsZero() || tx.Jsonrpc == "" ||
		tx.Method == "" || tx.Function == "" || tx.Args == nil {
		return ErrTxHasMissingProperties
	}
	return nil
}

func convertTxList(txList []command.Tx) []blockchain.Transaction {
	defaultTxList := make([]blockchain.Transaction, 0)

	for _, tx := range txList {
		defaultTx := convertTx(tx)
		defaultTxList = append(defaultTxList, defaultTx)
	}

	return defaultTxList
}

func convertTx(tx command.Tx) blockchain.Transaction {
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
