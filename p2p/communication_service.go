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

package p2p

import (
	"errors"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
)

var ErrEmptyConnectionId = errors.New("empty connection ")

type CallOnlyClient interface {
	Call(queue string, params interface{}, callback interface{}) error
}

type CommunicationService struct {
	client CallOnlyClient
}

func NewCommunicationService(client CallOnlyClient) *CommunicationService {

	return &CommunicationService{
		client: client,
	}
}

//dial to specific node
func (cs *CommunicationService) Dial(ipAddress string) error {

	c := command.CreateConnection{
		Address: ipAddress,
	}

	cs.client.Call("connection.create", c, func() {})

	return nil
}

//deliver peer leader table to specific peer
func (cs *CommunicationService) DeliverPLTable(connectionId string, peerLeaderTable PLTable) error {

	if connectionId == "" {
		return ErrEmptyConnectionId
	}

	if len(peerLeaderTable.PeerTable) == 0 {
		return ErrEmptyPeerTable
	}

	//create peer table message
	peerLeaderTableMessage := PLTableMessage{
		PLTable: peerLeaderTable,
	}

	grpcDeliverCommand, err := common.CreateGrpcDeliverCommand("PLTableDeliverProtocol", peerLeaderTableMessage)

	if err != nil {
		return err
	}

	grpcDeliverCommand.RecipientList = append(grpcDeliverCommand.RecipientList, connectionId)

	return cs.client.Call("message.deliver", grpcDeliverCommand, func() {})
}
