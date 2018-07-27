package blockchain_test

import (
	"testing"

	"encoding/hex"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/it-chain/engine/core/eventstore"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
)

func TestCommitBlock(t *testing.T) {

	//given
	tests := map[string]struct {
		input *blockchain.DefaultBlock
		err   error
	}{
		"success": {
			input: &blockchain.DefaultBlock{
				Seal:  []byte("Seal"),
				State: blockchain.Created,
			},
			err: nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		repo := mock.EventRepository{}

		repo.SaveFunc = func(aggregateID string, events ...midgard.Event) error {
			blockID := hex.EncodeToString([]byte("Seal"))
			assert.Equal(t, blockID, events[0].GetID())
			assert.Equal(t, 1, len(events))

			//when
			test.input.On(events[0])

			//then
			assert.Equal(t, blockchain.Committed, test.input.State)

			return nil
		}

		repo.CloseFunc = func() {}

		eventstore.InitForMock(repo)
		defer eventstore.Close()

		//when
		err := blockchain.CommitBlock(test.input)

		//then
		assert.Equal(t, test.err, err)

	}

}
