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

	"github.com/it-chain/midgard"
)

/*
 * consensus
 */

// Publish part

type PrepareMsgAdded struct {
	midgard.EventModel
	SenderId  string
	BlockHash []byte
}

type CommitMsgAdded struct {
	midgard.EventModel
	SenderId string
}

type ConsensusCreated struct {
	midgard.EventModel
	ConsensusId     string
	Representatives []*string
	Seal            []byte
	Body            []byte
	CurrentState    string
}

// Preprepare msg를 보냈을 때
type ConsensusPrePrepared struct {
	midgard.EventModel
}

// Prepare msg를 보냈을 때
type ConsensusPrepared struct {
	midgard.EventModel
}

// Commit msg를 보냈을 때
type ConsensusCommitted struct {
	midgard.EventModel
}

// block 저장이 끝나 state가 idle이 될 때
type ConsensusFinished struct {
	midgard.EventModel
}

// block이 저장되었을 때
type BlockSaved struct {
	midgard.EventModel
}

/*
 * grpc-gateway
 */

// ivm meta 생성
type MetaCreated struct {
	ICodeID        string
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        struct{}
}

// ivm meta deleted
type MetaDeleted struct {
	ICodeID string
}

// ivm meta status changed
type MetaStatusChanged struct {
	midgard.EventModel
	Status int
}

/*
 * blockChain
 */

// event when block is committed to event store
type BlockCommitted struct {
	midgard.EventModel
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []Tx
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   []byte
	State     string
}

// event when block is staged to event store
type BlockStaged struct {
	midgard.EventModel
	State string
}

//event when block is created in event store
type BlockCreated struct {
	midgard.EventModel
	Seal      []byte
	PrevSeal  []byte
	Height    uint64
	TxList    []Tx
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   []byte
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

type SyncStart struct {
	midgard.EventModel
}

type SyncDone struct {
	midgard.EventModel
}

/*
 * txpool
 */

// transaction created event
type TxCreated struct {
	midgard.EventModel
	ICodeID   string
	PeerID    string
	TimeStamp time.Time
	Jsonrpc   string
	Function  string
	Args      []string
	Signature []byte
}

// when block committed check transaction and delete
type TxDeleted struct {
	midgard.EventModel
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
	LeaderId string
}

//Connection

// connection 생성
type ConnectionCreated struct {
	ConnectionID string
	Address      string
}

// connection close
type ConnectionClosed struct {
	ConnectionId string
}
