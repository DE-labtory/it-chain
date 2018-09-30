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

package pbft_test

import (
	"testing"

	"time"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/stretchr/testify/assert"
)

func TestNewProposalPool(t *testing.T) {
	proposalPool := pbft.NewProposalPool()

	assert.NotNil(t, proposalPool)
}

func TestProposalPool_Pop(t *testing.T) {
	// given
	proposalPool := pbft.NewProposalPool()

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

	// when
	testProposal := proposalPool.Pop()

	// then
	assert.Equal(t, expectedCommand.Seal, testProposal.Seal)
}
