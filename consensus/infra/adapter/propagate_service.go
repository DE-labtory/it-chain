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
	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

type Publish func(exchange string, topic string, data interface{}) (err error)

type PropagateService struct {
	publish         Publish
	representatives []*consensus.Representative
}

func NewPropagateService(publish Publish, representatives []*consensus.Representative) *PropagateService {
	return &PropagateService{
		publish:         publish,
		representatives: representatives,
	}
}

func (ps PropagateService) BroadcastPrePrepareMsg(msg consensus.PrePrepareMsg) error {
	if msg.ConsensusId.Id == "" {
		return errors.New("Consensus ID is empty")
	}

	if msg.ProposedBlock.Body == nil {
		return errors.New("Block is empty")
	}

	SerializedMsg, err := common.Serialize(msg)

	if err != nil {
		return err
	}

	if err = ps.broadcastMsg(SerializedMsg, "SendPrePrepareMsgProtocol"); err != nil {
		return err
	}

	return nil
}

func (ps PropagateService) BroadcastPrepareMsg(msg consensus.PrepareMsg) error {
	if msg.ConsensusId.Id == "" {
		return errors.New("Consensus ID is empty")
	}

	if msg.BlockHash == nil {
		return errors.New("Block hash is empty")
	}

	SerializedMsg, err := common.Serialize(msg)

	if err != nil {
		return err
	}

	if err = ps.broadcastMsg(SerializedMsg, "SendPrepareMsgProtocol"); err != nil {
		return err
	}

	return nil
}

func (ps PropagateService) BroadcastCommitMsg(msg consensus.CommitMsg) error {
	if msg.ConsensusId.Id == "" {
		return errors.New("Consensus ID is empty")
	}

	SerializedMsg, err := common.Serialize(msg)

	if err != nil {
		return err
	}

	if err = ps.broadcastMsg(SerializedMsg, "SendCommitMsgProtocol"); err != nil {
		return err
	}

	return nil
}

func (ps PropagateService) broadcastMsg(SerializedMsg []byte, protocol string) error {
	if SerializedMsg == nil {
		return errors.New("Message is empty")
	}

	command, err := createDeliverGrpcCommand(protocol, SerializedMsg)

	if err != nil {
		return err
	}

	for _, r := range ps.representatives {
		command.RecipientList = append(command.RecipientList, r.GetID())
	}

	return ps.publish("Command", "message.broadcast", command)
}

func createDeliverGrpcCommand(protocol string, body interface{}) (command.DeliverGrpc, error) {
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
	}, nil
}
