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

	"github.com/it-chain/engine/grpc_gateway"
	"github.com/it-chain/engine/ivm"
)

/*
 * Consensus - pbft
 */

// Blockchain이 consensus를 요청하는 command
type StartConsensus struct {
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []Tx
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   string
	State     string
}

/*
 * grpc-gateway
 */

//Connection 생성 command
type CreateConnection struct {
	Address string
}

type JoinNetwork struct {
	Address string
}

//Connection close command
type CloseConnection struct {
	Address string
}

type GetConnectionList struct {
	ConnectionList []grpc_gateway.Connection
}

//다른 Peer에게 Message전송 command
type DeliverGrpc struct {
	MessageId     string
	RecipientList []string
	Body          []byte
	Protocol      string
}

//다른 Peer에게 Message수신 command
type ReceiveGrpc struct {
	MessageId    string
	Body         []byte
	ConnectionID string
	Protocol     string
}

type MyPeer struct {
	IpAddress string
	PeerId    string
}

/*
 * ivm
 */
type ExecuteICode struct {
	ICodeId  string
	Function string
	Args     []string
	Method   string
}

type Deploy struct {
	Url      string
	SshPath  string
	SshRaw   []byte
	Password string
}
type UnDeploy struct {
	ICodeId string
}

type GetICodeList struct {
}

type ICodeList struct {
	ICodes []ivm.ICode
}

/*
 * blockchain
 */

//Icode에게 block 내 TxList 실행 command
type ExecuteBlock struct {
	BlockId   string
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []Tx
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   string
	State     string
}

type ExecuteTransaction struct {
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

type ReturnBlockResult struct {
	BlockId      string
	TxResultList []TxResult
}

type TxResult struct {
	TxId    string
	Data    map[string]string
	Success bool
}

// Blockchain에게 block 생성 command
type ProposeBlock struct {
	BlockId string
	TxList  []Tx
}

type Tx struct {
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

/*
 * txpool
 */

type CreateTransaction struct {
	TransactionId string
	Jsonrpc       string
	Method        string
	ICodeID       string
	Function      string
	Args          []string
	Signature     []byte
}
