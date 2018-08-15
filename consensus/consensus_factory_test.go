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

package consensus_test

import (
	"testing"

	"github.com/it-chain/engine/consensus"
	"github.com/stretchr/testify/assert"
)

func TestCreateConsensus(t *testing.T) {
	// given
	p := make([]consensus.MemberId, 0)
	l := consensus.MemberId("leader")
	m := consensus.MemberId("member")
	b := consensus.ProposedBlock{
		Seal: make([]byte, 0),
		Body: make([]byte, 0),
	}

	// when
	c, err := consensus.CreateConsensus(p, b)

	// then
	assert.Error(t, err)

	// when
	p = append(p, l)
	p = append(p, m)

	c, err = consensus.CreateConsensus(p, b)

	// then
	assert.NoError(t, err)
	assert.Equal(t, 2, len(c.Representatives))
	assert.Equal(t, b.Seal, c.Block.Seal)
	assert.Equal(t, b.Body, c.Block.Body)
}

func TestConstructConsensus(t *testing.T) {
	// given
	l := consensus.NewRepresentative("leader")
	m := consensus.NewRepresentative("member")

	r := make([]*consensus.Representative, 0)
	r = append(r, l, m)

	msg := consensus.PrePrepareMsg{
		ConsensusId:    consensus.NewConsensusId("consensusID"),
		SenderId:       "me",
		Representative: r,
		ProposedBlock: consensus.ProposedBlock{
			Seal: make([]byte, 0),
			Body: make([]byte, 0),
		},
	}

	// when
	c, err := consensus.ConstructConsensus(msg)

	// then
	assert.NoError(t, err)
	assert.Equal(t, "consensusID", c.ConsensusID.Id)
	assert.Equal(t, consensus.IDLE_STATE, c.CurrentState)
	assert.Equal(t, 2, len(c.Representatives))
}
