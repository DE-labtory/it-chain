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
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/iLogger"
)

type transactionApiForSave interface {
	SaveTransactions(transactions []txpool.Transaction) error
}

type GrpcMessageHandler struct {
	transactionApi transactionApiForSave
}

func NewGrpcMessageHandler(transactionApi transactionApiForSave) *GrpcMessageHandler {
	return &GrpcMessageHandler{
		transactionApi: transactionApi,
	}
}

func (g GrpcMessageHandler) HandleMessageReceiveCommand(command command.ReceiveGrpc) {

	iLogger.Debug(nil, "[Txpool] Received grpc command")

	protocol := command.Protocol
	body := command.Body

	switch protocol {
	case txpool.SendTransactionsToLeader:

		transactionList := []txpool.Transaction{}

		if err := common.Deserialize(body, &transactionList); err != nil {
			iLogger.Errorf(nil, "[Txpool] Fail to deserialize grpcMessage - Err: [%s]", err.Error())
		}

		iLogger.Infof(nil, "[Txpool] Leader received transactions - Length: [%d]", len(transactionList))

		g.transactionApi.SaveTransactions(transactionList)

	}

}
