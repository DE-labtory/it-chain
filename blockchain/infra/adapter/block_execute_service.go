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
	"fmt"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type Publisher func(topic string, data interface{}) (err error) //해당 publish함수는 midgard 에서 의존성 주입을 받기 위해 interface로 작성한다.
//모든 의존성 주입은 컴포넌트.go 에서 이루어짐

type BlockExecuteService struct {
	client rpc.Client
}

func NewBlockExecuteService(client rpc.Client) *BlockExecuteService {
	return &BlockExecuteService{
		client: client,
	}
}

func (s *BlockExecuteService) ExecuteBlock(block blockchain.Block) error {

	deliverCommand, err := createBlockExecuteCommand(block)

	if err != nil {
		return err
	}

	err = s.client.Call("block.execute", deliverCommand, func(result command.ReturnBlockResult, err rpc.Error) {

		if !err.IsNil() {
			logger.Fatal(nil, err.Message)
			return
		}
		logger.Info(&logger.Fields{"": "BLOCKCHAIN"}, fmt.Sprintf("block executed %v", result))
	})

	return err
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
