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

package pbft

import (
	"testing"

	"time"

	"github.com/it-chain/engine/common/command"
	"github.com/stretchr/testify/assert"
)

func TestProposalPool_Save(t *testing.T) {
	// given
	proposalPool := NewProposalPool()

	assert.Equal(t, 0, len(proposalPool.proposals))

	expectedCommand := command.StartConsensus{
		Seal:      []byte{'s', 'e', 'a', 'l'},
		PrevSeal:  []byte{'p', 'r', 'e', 'v'},
		Height:    0,
		TxList:    make([]command.Tx, 0),
		TxSeal:    make([][]byte, 0),
		Timestamp: time.Time{},
		Creator:   "creator",
		State:     "state",
	}

	// when
	proposalPool.Save(expectedCommand)

	// then
	assert.Equal(t, 1, len(proposalPool.proposals))
	assert.Equal(t, expectedCommand.Seal, proposalPool.proposals[0].Seal)
}

func TestProposalPool_Pop(t *testing.T) {
	// given
	proposalPool := NewProposalPool()
	assert.Equal(t, 0, len(proposalPool.proposals))

	expectedCommand := command.StartConsensus{
		Seal:      []byte{'s', 'e', 'a', 'l'},
		PrevSeal:  []byte{'p', 'r', 'e', 'v'},
		Height:    0,
		TxList:    make([]command.Tx, 0),
		TxSeal:    make([][]byte, 0),
		Timestamp: time.Time{},
		Creator:   "creator",
		State:     "state",
	}

	proposalPool.Save(expectedCommand)
	assert.Equal(t, 1, len(proposalPool.proposals))

	// when
	testProposal := proposalPool.Pop()

	// then
	assert.Equal(t, 0, len(proposalPool.proposals))
	assert.Equal(t, expectedCommand.Seal, testProposal.Seal)
}

func TestProposalPool_RemoveAllMsg(t *testing.T) {
	// given
	proposalPool := NewProposalPool()
	assert.Equal(t, 0, len(proposalPool.proposals))

	expectedCommand1 := command.StartConsensus{
		Seal:      []byte{'S', 'E', 'A', 'L'},
		PrevSeal:  []byte{'p', 'r', 'e', 'v'},
		Height:    0,
		TxList:    make([]command.Tx, 0),
		TxSeal:    make([][]byte, 0),
		Timestamp: time.Time{},
		Creator:   "creator",
		State:     "state",
	}

	proposalPool.Save(expectedCommand1)
	assert.Equal(t, 1, len(proposalPool.proposals))

	expectedCommand2 := command.StartConsensus{
		Seal:      []byte{'s', 'e', 'a', 'l'},
		PrevSeal:  []byte{'S', 'E', 'A', 'L'},
		Height:    0,
		TxList:    make([]command.Tx, 0),
		TxSeal:    make([][]byte, 0),
		Timestamp: time.Time{},
		Creator:   "creator",
		State:     "state",
	}

	proposalPool.Save(expectedCommand2)
	assert.Equal(t, 2, len(proposalPool.proposals))

	// when
	proposalPool.RemoveAllMsg()

	// then
	assert.Equal(t, 0, len(proposalPool.proposals))
}
