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

	"sync"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/txpool"
	"github.com/rs/xid"
)

type BlockProposalService struct {
	client           rpc.Client // midgard.client
	engineMode       string
	txpoolRepository txpool.TransactionRepository
	sync.RWMutex
}

func NewBlockProposalService(client rpc.Client, txpoolRepository txpool.TransactionRepository, engineMode string) *BlockProposalService {
	return &BlockProposalService{
		client:           client,
		engineMode:       engineMode,
		RWMutex:          sync.RWMutex{},
		txpoolRepository: txpoolRepository,
	}
}

// todo do not delete transaction immediately
// todo transaction will be deleted when block are committed
func (b BlockProposalService) ProposeBlock() error {

	b.Lock()
	defer b.Unlock()

	// todo transaction size, number of tx
	transactions, err := b.txpoolRepository.FindAll()

	if err != nil {
		return err
	}

	if len(transactions) == 0 {
		return nil
	}

	if b.engineMode == "solo" {
		//propose transaction when solo mode
		if err := b.sendBlockProposal(transactions); err != nil {
			return err
		}

		for _, tx := range transactions {
			b.txpoolRepository.Remove(tx.ID)
		}

		return nil
	}

	return nil
}

func (b BlockProposalService) sendBlockProposal(transactions []txpool.Transaction) error {

	if len(transactions) == 0 {
		return errors.New("Empty transaction list proposed")
	}

	proposeCommand := command.ProposeBlock{
		BlockId: xid.New().String(),
		TxList:  make([]command.Tx, 0),
	}

	for _, tx := range transactions {
		proposeCommand.TxList = append(proposeCommand.TxList, command.Tx{
			ID:        tx.ID,
			PeerID:    tx.PeerID,
			ICodeID:   tx.ICodeID,
			TimeStamp: tx.TimeStamp,
			Jsonrpc:   tx.Jsonrpc,
			Function:  tx.Function,
			Args:      tx.Args,
			Signature: tx.Signature,
		})
	}

	err := b.client.Call("block.propose", proposeCommand, func(_ struct{}, err rpc.Error) {

		if !err.IsNil() {
			logger.Fatal(nil, err.Message)
			return
		}

		logger.Infof(nil, "[Txpool] Block has proposed")
	})

	return err
}
