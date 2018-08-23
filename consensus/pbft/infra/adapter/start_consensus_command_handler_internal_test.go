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

package adapter

import (
	"testing"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/stretchr/testify/assert"
)

func TestStartConsensusCommandHandler_extractProposedBlock(t *testing.T) {
	// given
	expectedSeal := []byte{'s', 'e', 'a', 'l'}
	expectedTxList := []command.Tx{}

	for i := 0; i < 5; i++ {
		expectedTxList = append(expectedTxList, command.Tx{
			ID: string(i),
		})
	}

	// when
	testBlock, err := extractProposedBlock(expectedSeal, expectedTxList)
	assert.NoError(t, err)

	expectedBody, err := common.Serialize(expectedTxList)
	assert.NoError(t, err)

	// then
	assert.Equal(t, expectedSeal, testBlock.Seal)
	assert.Equal(t, expectedBody, testBlock.Body)

	// given
	expectedSeal = nil

	// when
	testBlock, err = extractProposedBlock(expectedSeal, expectedTxList)

	// then
	assert.Equal(t, BlockSealIsNilError, err)
}
