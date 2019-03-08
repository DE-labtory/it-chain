/*
 * Copyright 2018 DE-labtory
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

package pbft

import (
	"errors"

	"github.com/DE-labtory/engine/common"
)

var ErrStateIdEmpty = errors.New("State ID is empty")
var ErrEmptyBlock = errors.New("Block is empty")
var ErrEmptyBlockHash = errors.New("Block hash is empty")
var ErrEmptyMsg = errors.New("Message is empty")

type PropagateService struct {
	eventService common.EventService
}

func NewPropagateService(eventService common.EventService) *PropagateService {
	return &PropagateService{
		eventService: eventService,
	}
}

func (ps PropagateService) BroadcastProposeMsg(msg ProposeMsg, representatives []Representative) error {

	if msg.StateID.ID == "" {
		return ErrStateIdEmpty
	}

	if msg.ProposedBlock.Body == nil {
		return ErrEmptyBlock
	}

	if err := ps.broadcastMsg(msg, "ProposeMsgProtocol", representatives); err != nil {
		return err
	}

	return nil
}

func (ps PropagateService) BroadcastPrevoteMsg(msg PrevoteMsg, representatives []Representative) error {
	if msg.StateID.ID == "" {
		return ErrStateIdEmpty
	}

	if msg.BlockHash == nil {
		return ErrEmptyBlockHash
	}

	if err := ps.broadcastMsg(msg, "PrevoteMsgProtocol", representatives); err != nil {
		return err
	}

	return nil
}

func (ps PropagateService) BroadcastPreCommitMsg(msg PreCommitMsg, representatives []Representative) error {
	if msg.StateID.ID == "" {
		return ErrStateIdEmpty
	}

	if err := ps.broadcastMsg(msg, "PreCommitMsgProtocol", representatives); err != nil {
		return err
	}

	return nil
}

func (ps PropagateService) broadcastMsg(msg interface{}, protocol string, representatives []Representative) error {

	grpcCommand, err := common.CreateGrpcDeliverCommand(protocol, msg)

	if err != nil {
		return err
	}

	for _, r := range representatives {
		grpcCommand.RecipientList = append(grpcCommand.RecipientList, r.GetID())
	}

	return ps.eventService.Publish("message.deliver", grpcCommand)
}
