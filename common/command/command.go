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

package command

import "github.com/it-chain/midgard"

/*
 * consensus
 */

type CreateBlockCommand struct {
	midgard.CommandModel
	Seal []byte
	Body []byte
}

type SendPrePrepareMsgCommand struct {
	midgard.CommandModel
	ConsensusId     string
	SenderId        string
	Representatives []*string
	Seal            []byte
	Body            []byte
}

type SendPrepareMsgCommand struct {
	midgard.CommandModel
	ConsensusId string
	SenderId    string
	BlockHash   []byte
}

type SendCommitMsgCommand struct {
	midgard.CommandModel
	ConsensusId string
	SenderId    string
}

type SendGrpcMsgCommand struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}

/*
 * grpc-gateway
 */

//Connection 생성 command
type ConnectionCreate struct {
	midgard.CommandModel
	Address string
}

//Connection close command
type ConnectionClose struct {
	midgard.CommandModel
}

//다른 Peer에게 Message전송 command
type GrpcDeliver struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}

//다른 Peer에게 Message수신 command
type GrpcReceive struct {
	midgard.CommandModel
	Body         []byte
	ConnectionID string
	Protocol     string
}
