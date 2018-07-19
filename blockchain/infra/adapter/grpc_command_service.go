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
	"github.com/it-chain/engine/common"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

// ToDo: 구현.(gitId:junk-sound)
type Publish func(exchange string, topic string, data interface{}) (err error)

type GrpcCommandService struct {
	publish Publish // midgard.client.Publish
}

func NewGrpcCommandService(publish Publish) *GrpcCommandService {
	return &GrpcCommandService{
		publish: publish,
	}
}

func (gcs *GrpcCommandService) RequestBlock(peerId blockchain.PeerId, height uint64) error {
	if peerId.Id == "" {
		return ErrEmptyNodeId
	}

	body := height

	deliverCommand, err := createGrpcDeliverCommand("BlockRequestProtocol", body)
	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, peerId.ToString())

	return gcs.publish("Command", "message.deliver", deliverCommand)
}

func (gcs *GrpcCommandService) ResponseBlock(peerId blockchain.PeerId, block blockchain.Block) error {
	if peerId.Id == "" {
		return ErrEmptyNodeId
	}

	if block.GetSeal() == nil {
		return ErrEmptyBlockSeal
	}

	body := block

	deliverCommand, err := createGrpcDeliverCommand("BlockResponseProtocol", body)
	if err != nil {
		return err
	}

	deliverCommand.Recipients = append(deliverCommand.Recipients, peerId.ToString())

	return gcs.publish("Command", "message.deliver", deliverCommand)
}

// TODO "SyncCheckResponseProtocol"을 통해서 last block을 전달한다.
func (gcs *GrpcCommandService) SyncCheckResponse(block blockchain.Block) error {
	return nil
}

func createGrpcDeliverCommand(protocol string, body interface{}) (blockchain.GrpcDeliverCommand, error) {

	data, err := common.Serialize(body)
	if err != nil {
		return blockchain.GrpcDeliverCommand{}, err
	}

	return blockchain.GrpcDeliverCommand{
		CommandModel: midgard.CommandModel{
			ID: xid.New().String(),
		},
		Recipients: make([]string, 0),
		Body:       data,
		Protocol:   protocol,
	}, err
}
