package adapter

import (
	"errors"
	"testing"

	"github.com/it-chain/engine/consensus"
	"github.com/stretchr/testify/assert"
)

func TestConfirmService_ConfirmBlock(t *testing.T) {
	tests := map[string]struct {
		input struct {
			block consensus.ProposedBlock
		}
		err error
	}{
		"success": {
			input: struct {
				block consensus.ProposedBlock
			}{
				block: consensus.ProposedBlock{
					Seal: make([]byte, 0),
					Body: make([]byte, 0),
				},
			},
			err: nil,
		},
		"block seal empty test": {
			input: struct {
				block consensus.ProposedBlock
			}{
				block: consensus.ProposedBlock{
					Seal: nil,
					Body: make([]byte, 0),
				},
			},
			err: errors.New("Block hash is nil"),
		},
		"block body empty test": {
			input: struct {
				block consensus.ProposedBlock
			}{
				block: consensus.ProposedBlock{
					Seal: make([]byte, 0),
					Body: nil,
				},
			},
			err: errors.New("There is no block"),
		},
	}

	publish := func(exchange string, topic string, data interface{}) (e error) {
		assert.Equal(t, "Command", exchange)
		assert.Equal(t, "block.create", topic)

		return nil
	}

	confirmService := NewConfirmService(publish)

	for testName, test := range tests {
		t.Logf("running test case [%s]", testName)

		err := confirmService.ConfirmBlock(test.input.block)

		assert.Equal(t, test.err, err)
	}
}
