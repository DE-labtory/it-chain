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

package adapter

import (
	"testing"
	"time"

	"github.com/DE-labtory/engine/common"
	"github.com/DE-labtory/engine/common/command"
	"github.com/stretchr/testify/assert"
)

func TestStartConsensusCommandHandler_extractProposedBlock(t *testing.T) {
	// given
	expectedSeal := []byte{'s', 'e', 'a', 'l'}
	expectedCommand := command.StartConsensus{
		Seal:      expectedSeal,
		PrevSeal:  []byte{'p', 'r', 'e', 'v'},
		Height:    0,
		TxList:    make([]command.Tx, 0),
		TxSeal:    make([][]byte, 0),
		Timestamp: time.Time{},
		Creator:   "creator",
		State:     "state",
	}

	// when
	testBlock, err := extractProposedBlock(expectedCommand)
	assert.NoError(t, err)

	expectedBody, err := common.Serialize(expectedCommand)
	assert.NoError(t, err)

	// then
	assert.Equal(t, expectedSeal, testBlock.Seal)
	assert.Equal(t, expectedBody, testBlock.Body)

	// given
	expectedSeal = nil
	expectedCommand.Seal = expectedSeal

	// when
	testBlock, err = extractProposedBlock(expectedCommand)

	// then
	assert.Equal(t, BlockSealIsNilError, err)
}
