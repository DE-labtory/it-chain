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

package consensus

import (
	"github.com/rs/xid"
)

// leader
func CreateConsensus(parliament []MemberId, block ProposedBlock) (*Consensus, error) {
	representatives, err := Elect(parliament)
	if err != nil {
		return &Consensus{}, err
	}

	newConsensus := Consensus{
		ConsensusID:     NewConsensusId(xid.New().String()),
		Representatives: representatives,
		Block:           block,
		CurrentState:    IDLE_STATE,
		PrepareMsgPool:  NewPrepareMsgPool(),
		CommitMsgPool:   NewCommitMsgPool(),
	}

	return &newConsensus, nil
}

// member
func ConstructConsensus(msg PrePrepareMsg) (*Consensus, error) {
	newConsensus := &Consensus{
		ConsensusID:     msg.ConsensusId,
		Representatives: msg.Representative,
		Block:           msg.ProposedBlock,
		CurrentState:    IDLE_STATE,
		PrepareMsgPool:  NewPrepareMsgPool(),
		CommitMsgPool:   NewCommitMsgPool(),
	}

	return newConsensus, nil
}
