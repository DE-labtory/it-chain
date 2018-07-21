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
	"github.com/it-chain/midgard"
	"github.com/pkg/errors"
	"github.com/rs/xid"
)

var ErrBlockTypeCasting = errors.New("Error in type casting block")

type BlockExecuteService struct {
	publisher Publisher
}

func NewBlockExecuteService(publisher Publisher) *BlockExecuteService {
	return &BlockExecuteService{
		publisher: publisher,
	}
}

func (s *BlockExecuteService) ExecuteBlock(block blockchain.Block) error {
	command, err := createBlockExecuteCommand(block)
	if err != nil {
		return err
	}

	return s.publisher("block.execute", command)
}

func createBlockExecuteCommand(block blockchain.Block) (command.ExecuteBlock, error) {
	defaultBlock, ok := block.(*blockchain.DefaultBlock)
	if !ok {
		return command.ExecuteBlock{}, ErrBlockTypeCasting
	}

	return command.ExecuteBlock{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Seal:      (*defaultBlock).Seal,
		PrevSeal:  (*defaultBlock).PrevSeal,
		Height:    (*defaultBlock).Height,
		TxList:    convertToExecuteBlockTxList((*defaultBlock).TxList),
		TxSeal:    (*defaultBlock).TxSeal,
		Timestamp: (*defaultBlock).Timestamp,
		Creator:   (*defaultBlock).Creator,
		State:     (*defaultBlock).State,
	}, nil
}

func convertToExecuteBlockTxList(txList []*blockchain.DefaultTransaction) []command.Tx {
	executeBlockTxList := make([]command.Tx, 0)

	for _, tx := range txList {
		executeBlockTx := convertToExecuteBlockTx(tx)
		executeBlockTxList = append(executeBlockTxList, executeBlockTx)
	}

	return executeBlockTxList
}

func convertToExecuteBlockTx(tx *blockchain.DefaultTransaction) command.Tx {
	return command.Tx{
		ID:        tx.ID,
		ICodeID:   tx.ICodeID,
		Status:    int(tx.Status),
		PeerID:    tx.PeerID,
		TimeStamp: tx.Timestamp,
		Jsonrpc:   tx.TxData.Jsonrpc,
		Method:    string(tx.TxData.Method),
		Args:      tx.TxData.Params.Args,
		Function:  tx.TxData.Params.Function,
		Signature: tx.Signature,
	}
}
