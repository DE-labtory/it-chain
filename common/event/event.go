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

// Consume part
type LeaderChanged struct {
	midgard.EventModel
	LeaderId string
}

type MemberJoined struct {
	midgard.EventModel
	MemberId string
}

type MemberRemoved struct {
	midgard.EventModel
	MemberId string
}

// block이 저장되었을 때
type BlockSaved struct {
	midgard.EventModel
}

/*
 * grpc-gateway
 */

// connection 생성
type ConnectionCreated struct {
	midgard.EventModel
	Address string
}

// connection close
type ConnectionClosed struct {
	midgard.EventModel
}

// icode meta 생성
type MetaCreated struct {
	midgard.EventModel
	RepositoryName string
	GitUrl         string
	Path           string
	CommitHash     string
	Version        struct{}
}

// icode meta deleted
type MetaDeleted struct {
	midgard.EventModel
}

// icode meta status changed
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
	State string
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
	TxList    []byte
	TxSeal    [][]byte
	Timestamp time.Time
	Creator   []byte
	State     string
}

/*
 * txpool
 */

// transaction created event
type TxCreated struct {
	midgard.EventModel
	Status    int
	TimeStamp time.Time
	Jsonrpc   string
	Method    string
	ICodeID   string
	Function  string
	Args      []string
	Signature []byte
	PeerID    string
}

// when block committed check transaction and delete
type TxDeleted struct {
	midgard.EventModel
}
