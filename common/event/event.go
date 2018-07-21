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

import "github.com/it-chain/midgard"

/*
 * consensus
 */

// Publish part

type PrepareMsgAddedEvent struct {
	midgard.EventModel
	SenderId  string
	BlockHash []byte
}

type CommitMsgAddedEvent struct {
	midgard.EventModel
	SenderId string
}

type ConsensusCreatedEvent struct {
	midgard.EventModel
	ConsensusId     string
	Representatives []*string
	Seal            []byte
	Body            []byte
	CurrentState    string
}

// Preprepare msg를 보냈을 때
type ConsensusPrePreparedEvent struct {
	midgard.EventModel
}

// Prepare msg를 보냈을 때
type ConsensusPreparedEvent struct {
	midgard.EventModel
}

// Commit msg를 보냈을 때
type ConsensusCommittedEvent struct {
	midgard.EventModel
}

// block 저장이 끝나 state가 idle이 될 때
type ConsensusFinishedEvent struct {
	midgard.EventModel
}

// Consume part
type LeaderChangedEvent struct {
	midgard.EventModel
	LeaderId string
}

type MemberJoinedEvent struct {
	midgard.EventModel
	MemberId string
}

type MemberRemovedEvent struct {
	midgard.EventModel
	MemberId string
}

// block이 저장되었을 때
type BlockSavedEvent struct {
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
