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

package adapter_test

import (
	"testing"

	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/infra/adapter"
	"github.com/stretchr/testify/assert"
)

func TestConfirmService_ConfirmBlock(t *testing.T) {
	tests := map[string]struct {
		input struct {
			block pbft.ProposedBlock
		}
		err error
	}{
		"success": {
			input: struct {
				block pbft.ProposedBlock
			}{
				block: pbft.ProposedBlock{
					Seal: make([]byte, 0),
					Body: make([]byte, 0),
				},
			},
			err: nil,
		},
		"block seal empty test": {
			input: struct {
				block pbft.ProposedBlock
			}{
				block: pbft.ProposedBlock{
					Seal: nil,
					Body: make([]byte, 0),
				},
			},
			err: adapter.ErrBlockHashNil,
		},
		"block body empty test": {
			input: struct {
				block pbft.ProposedBlock
			}{
				block: pbft.ProposedBlock{
					Seal: make([]byte, 0),
					Body: nil,
				},
			},
			err: adapter.ErrNoBlock,
		},
	}

	publish := func(topic string, data interface{}) (e error) {
		assert.Equal(t, "block.confirm", topic)

		return nil
	}

	confirmService := adapter.NewConfirmService(publish)

	for testName, test := range tests {
		t.Logf("running test case [%s]", testName)

		err := confirmService.ConfirmBlock(test.input.block)

		assert.Equal(t, test.err, err)
	}
}
