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

// consensus가 끝났다는 event
// true면 저장, false면 저장안함
type ConsensusFinished struct {
	isConfirm bool
}

/*
 * grpc-gateway
 */

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
	midgard.EventModel
	IpAddress string
}

type PeerDeleted struct {
	midgard.EventModel
}

// handle leader received event
type LeaderUpdated struct {
	midgard.EventModel
}

type LeaderDelivered struct {
	midgard.EventModel
}

type LeaderDeleted struct {
	midgard.EventModel
}

//Connection

// connection 생성
type ConnectionCreated struct {
	midgard.EventModel
	Address string
}

// connection close
type ConnectionClosed struct {
	midgard.EventModel
}
