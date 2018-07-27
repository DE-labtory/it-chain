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
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/midgard"
	"github.com/rs/xid"
)

// leader
func CreateConsensus(parliament []MemberId, block ProposedBlock) (*Consensus, error) {
	representatives, err := Elect(parliament)
	if err != nil {
		return &Consensus{}, err
	}

	consensusID := NewConsensusId(xid.New().String())
	consensus := &Consensus{}
	consensusCreatedEvent := newConsensusCreatedEvent(consensusID, representatives, block)

	if err := OnAndSave(consensus, &consensusCreatedEvent); err != nil {
		return &Consensus{}, err
	}

	return consensus, nil
}

// member
func ConstructConsensus(msg PrePrepareMsg) (*Consensus, error) {
	consensus := &Consensus{}
	consensusCreatedEvent := newConsensusCreatedEvent(msg.ConsensusId, msg.Representative, msg.ProposedBlock)

	if err := OnAndSave(consensus, &consensusCreatedEvent); err != nil {
		return &Consensus{}, err
	}

	return consensus, nil
}

func newConsensusCreatedEvent(cID ConsensusId, r []*Representative, b ProposedBlock) event.ConsensusCreated {
	representativeStr := make([]*string, 0)

	for _, representative := range r {
		str := representative.GetID()
		representativeStr = append(representativeStr, &str)
	}

	return event.ConsensusCreated{
		EventModel: midgard.EventModel{
			ID: cID.Id,
		},
		ConsensusId:     cID.Id,
		Representatives: representativeStr,
		Seal:            b.Seal,
		Body:            b.Body,
		CurrentState:    string(IDLE_STATE),
	}
}
