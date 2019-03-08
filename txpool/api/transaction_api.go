/*
 * Copyright 2018 DE-labtory
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
	"github.com/DE-labtory/engine/txpool"
	"github.com/DE-labtory/iLogger"
)

type TransactionApi struct {
	nodeId                string
	transactionRepository txpool.TransactionRepository
	leaderRepository      txpool.LeaderRepository
	transferService       *txpool.TransferService
	blockProposalService  *txpool.BlockProposalService
}

func NewTransactionApi(nodeId string, transactionRepository txpool.TransactionRepository, leaderRepository txpool.LeaderRepository, transferService *txpool.TransferService, blockProposalService *txpool.BlockProposalService) *TransactionApi {
	return &TransactionApi{
		nodeId:                nodeId,
		transactionRepository: transactionRepository,
		leaderRepository:      leaderRepository,
		transferService:       transferService,
		blockProposalService:  blockProposalService,
	}
}

func (t TransactionApi) CreateTransaction(txData txpool.TxData) (txpool.Transaction, error) {

	transaction, err := txpool.CreateTransaction(t.nodeId, txData)

	if err != nil {
		iLogger.Errorf(nil, "[Txpool] Fail to create transaction - Err: [%s]", err)
		return txpool.Transaction{}, err
	}

	err = t.transactionRepository.Save(transaction)

	return transaction, err
}

func (t TransactionApi) SaveTransactions(transactions []txpool.Transaction) error {

	for _, tx := range transactions {

		if err := t.transactionRepository.Save(tx); err != nil {
			return err
		}
	}

	return nil
}

func (t TransactionApi) DeleteTransaction(id txpool.TransactionId) {

	t.transactionRepository.Remove(id)
}

func (t TransactionApi) ProposeBlock(engineMode string) error {

	switch engineMode {

	case "solo":
		return t.blockProposalService.ProposeBlock()

	case "pbft":
		if t.isLeader() {
			return t.blockProposalService.ProposeBlock()
		}

		return nil

	default:
		return nil
	}
}

func (t TransactionApi) SendLeaderTransaction(engineMode string) error {

	switch engineMode {

	case "pbft":
		if t.isLeader() {
			return nil
		}

		if !t.isLeaderExistent() {
			return nil
		}

		return t.transferService.SendLeaderTransactions()

	default:
		return nil
	}
}

func (t TransactionApi) isLeader() bool {

	leader := t.leaderRepository.Get()

	return txpool.IsLeader(t.nodeId, leader)
}

func (t TransactionApi) isLeaderExistent() bool {
	leader := t.leaderRepository.Get()

	if leader.Id == "" {
		return false
	}

	return true
}
