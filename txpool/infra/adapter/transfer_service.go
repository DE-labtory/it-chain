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

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

var ErrTxEmpty = errors.New("Empty transaction list proposed")

type Publisher func(topic string, data interface{}) (err error) //해당 publish함수는 midgard 에서 의존성 주입을 받기 위해 interface로 작성한다.
//모든 의존성 주입은 컴포넌트.go 에서 이루어짐

type TransferService struct {
	publisher Publisher // midgard.client
}

func NewTransferService(publisher Publisher) *TransferService {
	return &TransferService{
		publisher: publisher,
	}
}

func (ts TransferService) SendLeaderTransactions(transactions []txpool.Transaction, leader txpool.Leader) error {

	if len(transactions) == 0 {
		return ErrTxEmpty
	}

	deliverCommand, err := createGrpcDeliverCommand("SendLeaderTransactionsProtocol", transactions)

	if err != nil {
		return err
	}

	deliverCommand.RecipientList = append(deliverCommand.RecipientList, leader.LeaderId.ToString())

	return ts.publisher("message.deliver", deliverCommand)
}

func createGrpcDeliverCommand(protocol string, body interface{}) (command.DeliverGrpc, error) {

	data, err := common.Serialize(body)

	if err != nil {
		return command.DeliverGrpc{}, err
	}

	return command.DeliverGrpc{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		RecipientList: make([]string, 0),
		Body:          data,
		Protocol:      protocol,
	}, err
}
