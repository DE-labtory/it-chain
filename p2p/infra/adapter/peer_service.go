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
	"time"

	"github.com/it-chain/engine/grpc_gateway"
	"github.com/it-chain/engine/p2p"
)

type PeerService struct {
	publish Publish
}

func (ps *PeerService) Dial(ipAddress string) error {
	command := grpc_gateway.ConnectionCreateCommand{
		Address: ipAddress,
	}
	ps.publish("Command", "connection.create", command)
	return nil
}

//request leader information in p2p network to the node specified by peerId
func (ps *PeerService) RequestLeaderInfo(connectionId string) error {

	if connectionId == "" {
		return ErrEmptyPeerId
	}

	body := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	//message deliver command for delivering leader info
	deliverCommand, err := CreateGrpcDeliverCommand("LeaderInfoRequestMessage", body)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, connectionId)

	return ps.publish("Command", "message.deliver", deliverCommand)
}

// command message which requests node list of specific node
func (ps *PeerService) RequestPeerList(peerId p2p.PeerId) error {

	if peerId.Id == "" {
		return ErrEmptyPeerId
	}
	body := p2p.LeaderInfoRequestMessage{
		TimeUnix: time.Now().Unix(),
	}

	deliverCommand, err := CreateGrpcDeliverCommand("PeerListRequestMessage", body)

	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, peerId.ToString())

	return ps.publish("Command", "message.deliver", deliverCommand)
}
