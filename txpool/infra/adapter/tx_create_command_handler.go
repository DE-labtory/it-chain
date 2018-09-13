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
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/engine/txpool/api"
)

type TxCommandHandler struct {
	transactionApi *api.TransactionApi
}

func NewTxCommandHandler(transactionApi *api.TransactionApi) *TxCommandHandler {
	return &TxCommandHandler{
		transactionApi: transactionApi,
	}
}

func (t *TxCommandHandler) HandleTxCreateCommand(txCreateCommand command.CreateTransaction) (txpool.Transaction, rpc.Error) {

	txData := txpool.TxData{
		ICodeID:   txCreateCommand.ICodeID,
		Jsonrpc:   txCreateCommand.Jsonrpc,
		Function:  txCreateCommand.Function,
		Signature: txCreateCommand.Signature,
		Args:      txCreateCommand.Args,
	}

	tx, err := t.transactionApi.CreateTransaction(txData)

	if err != nil {
		return txpool.Transaction{}, rpc.Error{Message: err.Error()}
	}

	return tx, rpc.Error{}
}
