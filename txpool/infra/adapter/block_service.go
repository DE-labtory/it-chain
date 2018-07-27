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

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type BlockService struct {
	client rpc.Client // midgard.client
}

func NewBlockService(client rpc.Client) *BlockService {
	return &BlockService{
		client: client,
	}
}

func (b BlockService) ProposeBlock(transactions []txpool.Transaction) error {

	if len(transactions) == 0 {
		return errors.New("Empty transaction list proposed")
	}

	proposeCommand := command.ProposeBlock{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},

		TxList: make([]command.Tx, 0),
	}

	for _, tx := range transactions {
		proposeCommand.TxList = append(proposeCommand.TxList, command.Tx{
			ID:        tx.ID,
			Status:    int(tx.Status),
			PeerID:    tx.PeerID,
			TimeStamp: tx.TimeStamp,
			Jsonrpc:   tx.Jsonrpc,
			Method:    string(tx.Method),
			Function:  tx.Function,
			Args:      tx.Args,
			Signature: tx.Signature,
		})
	}

	err := b.client.Call("block.propose", proposeCommand, func() {

	})

	return err
}
