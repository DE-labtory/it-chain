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
	"github.com/it-chain/engine/p2p"
)

type CommunicationService struct {
	publish Publish
}

func NewCommunicationService(publish Publish) *CommunicationService {

	return &CommunicationService{
		publish: publish,
	}
}

func (cs *CommunicationService) Dial(ipAddress string) error {

	command := p2p.ConnectionCreateCommand{
		Address: ipAddress,
	}

	cs.publish("Command", "connection.create", command)

	return nil
}

func (cs *CommunicationService) DeliverPLTable(connectionId string, peerLeaderTable p2p.PLTable) error {

	if connectionId == "" {
		return ErrEmptyConnectionId
	}

	if len(peerLeaderTable.PeerTable) == 0 {
		return p2p.ErrEmptyPeerTable
	}

	//create peer table message
	peerLeaderTableMessage := p2p.PLTableMessage{
		PLTable: peerLeaderTable,
	}

	grpcDeliverCommand, err := CreateGrpcDeliverCommand("PLTableDeliverProtocol", peerLeaderTableMessage)

	if err != nil {
		return err
	}

	grpcDeliverCommand.Recipients = append(grpcDeliverCommand.Recipients, connectionId)

	return cs.publish("Command", "message.deliver", grpcDeliverCommand)
}
