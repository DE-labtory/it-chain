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
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/iLogger"
)

var DeserializingError = errors.New("Message deserializing is failed.")
var UndefinedProtocolError = errors.New("Undefined protocol error")

type StateMsgApi interface {
	HandleProposeMsg(msg pbft.ProposeMsg) error
	HandlePrevoteMsg(msg pbft.PrevoteMsg) error
	HandlePreCommitMsg(msg pbft.PreCommitMsg) error
}

type PbftMsgHandler struct {
	sApi StateMsgApi
}

func NewPbftMsgHandler(sApi StateMsgApi) *PbftMsgHandler {
	return &PbftMsgHandler{
		sApi: sApi,
	}
}

func (p *PbftMsgHandler) HandleGrpcMsgCommand(command command.ReceiveGrpc) error {
	protocol := command.Protocol
	body := command.Body

	iLogger.Infof(nil, "[PBFT] Received protocol - Protocol: [%s]", protocol)
	switch protocol {

	case "ProposeMsgProtocol":
		msg := pbft.ProposeMsg{}
		if err := common.Deserialize(body, &msg); err != nil {
			logger.Errorf(nil, "[PBFT] %s", DeserializingError.Error())
		}

		if err := p.sApi.HandleProposeMsg(msg); err != nil {
			logger.Errorf(nil, "[PBFT] %s", err.Error())
		}

	case "PrevoteMsgProtocol":
		msg := pbft.PrevoteMsg{}
		if err := common.Deserialize(body, &msg); err != nil {
			logger.Errorf(nil, "[PBFT] %s", DeserializingError.Error())
		}

		if err := p.sApi.HandlePrevoteMsg(msg); err != nil {
			logger.Errorf(nil, "[PBFT] %s", err.Error())
		}

	case "PreCommitMsgProtocol":
		msg := pbft.PreCommitMsg{}
		if err := common.Deserialize(body, &msg); err != nil {
			logger.Errorf(nil, "[PBFT] %s", DeserializingError.Error())
		}

		if err := p.sApi.HandlePreCommitMsg(msg); err != nil {
			logger.Errorf(nil, "[PBFT] %s", err.Error())
		}

	default:
		logger.Errorf(nil, "[PBFT] %s", UndefinedProtocolError.Error())
	}
	return nil
}
