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
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/consensus"
	"github.com/it-chain/engine/consensus/api"
)

type GrpcMsgCommandHandler struct {
	cApi api.ConsensusApi
}

func NewGrpcMsgCommandHandler(cApi api.ConsensusApi) GrpcMsgCommandHandler {
	return GrpcMsgCommandHandler{}
}

func (s *GrpcMsgCommandHandler) HandleGrpcMsgCommand(command command.ReceiveGrpc) {
	protocol := command.Protocol
	body := command.Body

	switch protocol {
	case "PrePrepareMsgProtocol":
		s.HandleSendPrePrepareMsg(body)
	case "PrepareMsgProtocol":
		s.HandleSendPrepareMsg(body)
	case "CommitMsgProtocol":
		s.HandleSendCommitMsg(body)
	default:
		logger.Error(nil, "GRPC protocol is not defined.")
	}
}

func (s *GrpcMsgCommandHandler) HandleSendPrePrepareMsg(body []byte) {
	var msg consensus.PrePrepareMsg
	common.Deserialize(body, msg)

	s.cApi.ReceivePrePrepareMsg(msg)
}

func (s *GrpcMsgCommandHandler) HandleSendPrepareMsg(body []byte) {
	var msg consensus.PrepareMsg
	common.Deserialize(body, msg)

	s.cApi.ReceivePrepareMsg(msg)
}

func (s *GrpcMsgCommandHandler) HandleSendCommitMsg(body []byte) {
	var msg consensus.CommitMsg
	common.Deserialize(body, msg)

	s.cApi.ReceiveCommitMsg(msg)
}
