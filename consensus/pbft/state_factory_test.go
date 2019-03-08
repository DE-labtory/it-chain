/*
 * Copyright 2018 DE-labtory
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

package pbft_test

import (
	"testing"

	"github.com/DE-labtory/engine/consensus/pbft"
	"github.com/stretchr/testify/assert"
)

func TestCreateConsensus(t *testing.T) {
	// given
	p := make([]pbft.Representative, 0)
	l := pbft.Representative{
		ID: "leader",
	}
	m := pbft.Representative{
		ID: "member",
	}
	b := pbft.ProposedBlock{
		Seal: make([]byte, 0),
		Body: make([]byte, 0),
	}

	// when
	c, err := pbft.NewState(p, b)

	// then
	assert.Error(t, err)

	// when
	p = append(p, l)
	p = append(p, m)

	c, err = pbft.NewState(p, b)

	// then
	assert.NoError(t, err)
	assert.Equal(t, 2, len(c.Representatives))
	assert.Equal(t, b.Seal, c.Block.Seal)
	assert.Equal(t, b.Body, c.Block.Body)
}

func TestConstructConsensus(t *testing.T) {
	// given
	l := pbft.NewRepresentative("leader")
	m := pbft.NewRepresentative("member")

	r := make([]pbft.Representative, 0)
	r = append(r, l, m)

	msg := pbft.ProposeMsg{
		StateID:        pbft.NewStateID("consensusID"),
		SenderID:       "me",
		Representative: r,
		ProposedBlock: pbft.ProposedBlock{
			Seal: make([]byte, 0),
			Body: make([]byte, 0),
		},
	}

	// when
	c := pbft.BuildState(msg)

	// then
	assert.Equal(t, "consensusID", c.StateID.ID)
	assert.Equal(t, pbft.IDLE_STAGE, c.CurrentStage)
	assert.Equal(t, 2, len(c.Representatives))
}
