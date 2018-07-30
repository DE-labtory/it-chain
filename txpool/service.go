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
	"log"
	"sync"
)

type TxpoolQueryService interface {
	FindUncommittedTransactions() ([]Transaction, error)
}

type TransferService interface {
	SendTransactionsToLeader(transactions []Transaction, leader Leader) error
}

type BlockService interface {
	ProposeBlock(transactions []Transaction) error
}

type BlockProposalService struct {
	engineMode         string
	txpoolQueryService TxpoolQueryService
	blockService       BlockService
	sync.RWMutex
}

func NewBlockProposalService(queryService TxpoolQueryService, blockService BlockService, engineMode string) *BlockProposalService {

	return &BlockProposalService{
		txpoolQueryService: queryService,
		blockService:       blockService,
		engineMode:         engineMode,
		RWMutex:            sync.RWMutex{},
	}
}

// todo do not delete transaction immediately
// todo transaction will be deleted when block are committed
func (b BlockProposalService) ProposeBlock() error {

	b.Lock()
	defer b.Unlock()

	// todo transaction size, number of tx
	transactions, err := b.txpoolQueryService.FindUncommittedTransactions()

	if err != nil {
		return err
	}

	if len(transactions) == 0 {
		return nil
	}

	if b.engineMode == "solo" {
		//propose transaction when solo mode
		err = b.blockService.ProposeBlock(transactions)

		if err != nil {
			return err
		}

		log.Printf("transactions are proposed [%v]", transactions)

		for _, tx := range transactions {
			DeleteTransaction(tx)
		}

		return nil
	}

	return nil
}

func filter(vs []Transaction, f func(Transaction) bool) []Transaction {
	vsf := make([]Transaction, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
