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

package event

import (
	"time"
)

/*
 * consensus
 */

// todo 블록 정보 가지고 있어야 함
// consensus가 끝났다는 event
// true면 블록 저장, false면 블록 저장 안함
type ConsensusFinished struct {
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

// ivm meta 생성
type ICodeCreated struct {
	ID             string
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        string
}

// ivm meta deleted
type ICodeDeleted struct {
	ICodeID string
}

/*
 * blockChain
 */

// event when block is committed to event store
type BlockCommitted struct {
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []Tx
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   string
	State     string
}

// event when block is staged to event store
type BlockStaged struct {
	BlockId string
	State   string
}

//event when block is created in event store
type BlockCreated struct {
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

type Tx struct {
	ID        string
	ICodeID   string
	PeerID    string
	TimeStamp time.Time
	Jsonrpc   string
	Function  string
	Args      []string
	Signature []byte
}

/*
 * txpool
 */

// transaction created event
type TxCreated struct {
	TransactionId string
	ICodeID       string
	PeerID        string
	TimeStamp     time.Time
	Jsonrpc       string
	Function      string
	Args          []string
	Signature     []byte
}

// when block committed check transaction and delete
type TxDeleted struct {
	TransactionId string
}

/*
 * p2p
 */

type PeerCreated struct {
	PeerId    string
	IpAddress string
}

type PeerDeleted struct {
	PeerId string
}

// handle leader received event
type LeaderUpdated struct {
	LeaderId string
}

type LeaderDelivered struct {
	LeaderId string
}

type LeaderDeleted struct {
}

//Connection

// connection 생성
type ConnectionCreated struct {
	ConnectionID       string
	GrpcGatewayAddress string
	ApiGatewayAddress  string
}

type ConnectionSaved struct {
	ConnectionID string
	Address      string
}

// connection close
type ConnectionClosed struct {
	ConnectionID string
}
