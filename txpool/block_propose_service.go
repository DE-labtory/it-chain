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
	"sync"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/iLogger"
)

type BlockProposalService struct {
	txpoolRepository TransactionRepository
	eventService     EventService
	sync.RWMutex
}

func NewBlockProposalService(txpoolRepository TransactionRepository, eventService EventService) *BlockProposalService {
	return &BlockProposalService{
		txpoolRepository: txpoolRepository,
		eventService:     eventService,
		RWMutex:          sync.RWMutex{},
	}
}

// todo do not delete transaction immediately
// todo transaction will be deleted when block are committed
func (b BlockProposalService) ProposeBlock() error {

	b.Lock()
	defer b.Unlock()

	// todo transaction size, number of tx
	transactions, err := b.txpoolRepository.FindAll()

	iLogger.Debugf(nil, "[Txpool] transaction number - tx: [%d]", len(transactions))
	if err != nil {
		return err
	}

	if len(transactions) == 0 {
		return nil
	}

	if err := b.sendBlockProposal(transactions); err != nil {
		return err
	}

	b.clearTransactions(transactions)

	return nil

}

func (b BlockProposalService) sendBlockProposal(transactions []Transaction) error {

	ProposeBlockEvent := createProposeBlockCommand(transactions)

	if err := b.eventService.Publish("block.propose", ProposeBlockEvent); err != nil {
		iLogger.Errorf(nil, "[Txpool] Fail to propose block - Err: [%s]", err)
		return err
	}

	iLogger.Info(nil, "[Txpool] Block has proposed")
	return nil
}

func (b BlockProposalService) clearTransactions(transactions []Transaction) {
	for _, tx := range transactions {
		b.txpoolRepository.Remove(tx.ID)
	}
}

func createProposeBlockCommand(transactions []Transaction) command.ProposeBlock {
	return command.ProposeBlock{
		TxList: convertTxListType(transactions),
	}
}

func convertTxListType(transactions []Transaction) []command.Tx {
	txList := make([]command.Tx, 0)
	for _, tx := range transactions {
		txList = append(txList, convertTxType(tx))
	}

	return txList
}

func convertTxType(tx Transaction) command.Tx {
	return command.Tx{
		ID:        tx.ID,
		ICodeID:   tx.ICodeID,
		PeerID:    tx.PeerID,
		TimeStamp: tx.TimeStamp,
		Jsonrpc:   tx.Jsonrpc,
		Function:  tx.Function,
		Args:      tx.Args,
		Signature: tx.Signature,
	}
}
