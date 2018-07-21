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

import (
	"time"

	"github.com/it-chain/midgard"
)

/*
 * consensus
 */

type ConfirmBlock struct {
	midgard.CommandModel
	Seal []byte
	Body []byte
}

type SendPrePrepareMsg struct {
	midgard.CommandModel
	ConsensusId     string
	SenderId        string
	Representatives []*string
	Seal            []byte
	Body            []byte
}

type SendPrepareMsg struct {
	midgard.CommandModel
	ConsensusId string
	SenderId    string
	BlockHash   []byte
}

type SendCommitMsg struct {
	midgard.CommandModel
	ConsensusId string
	SenderId    string
}

type SendGrpcMsg struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}

/*
 * grpc-gateway
 */

//Connection 생성 command
type CreateConnection struct {
	midgard.CommandModel
	Address string
}

//Connection close command
type CloseConnection struct {
	midgard.CommandModel
}

//다른 Peer에게 Message전송 command
type DeliverGrpc struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}

//다른 Peer에게 Message수신 command
type ReceiveGrpc struct {
	midgard.CommandModel
	Body         []byte
	ConnectionID string
	Protocol     string
}

/*
 * icode
 */
type ExecuteTransaction struct {
	midgard.CommandModel
	ICodeID   string
	Status    int
	PeerID    string
	TimeStamp time.Time
	Jsonrpc   string
	Method    string
	Function  string
	Args      []string
	Signature []byte
}

type Deploy struct {
	midgard.CommandModel
	Url     string
	SshPath string
}
type UnDeploy struct {
	midgard.CommandModel
}

type Query struct {
	midgard.CommandModel
	Function string
	Args     []string
}

/*
 * blockchain
 */

//Icode에게 block 내 TxList 실행 command
type ExecuteBlock struct {
	midgard.CommandModel
	Seal     []byte
	PrevSeal []byte
	Height   uint64
	TxList   []struct {
		ID        string
		ICodeID   string
		Status    int
		PeerID    string
		TimeStamp time.Time
		Jsonrpc   string
		Method    string
		Function  string
		Args      []string
		Signature []byte
	}
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   []byte
	State     string
}

type ReturnBlockResult struct {
	midgard.CommandModel
	TxResults []struct {
		TxId    string
		Data    map[string]string
		Success bool
	}
}

// Blockchain에게 block 생성 command
type ProposeBlock struct {
	midgard.CommandModel
	TxList []ProposeBlockTx
}

type ProposeBlockTx struct {
	ID        string
	Status    int
	PeerID    string
	TimeStamp time.Time
	Jsonrpc   string
	Method    string
	Type      int
	Function  string
	Args      []string
	Signature []byte
}
