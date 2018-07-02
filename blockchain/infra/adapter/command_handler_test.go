package adapter_test

import (
	"testing"
	"github.com/it-chain/midgard"
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/infra/adapter"
	"github.com/stretchr/testify/assert"
)

type MockBlockApi struct {
	AddBlockToPoolFunc func(block blockchain.Block)
}
func (api MockBlockApi) AddBlockToPool(block blockchain.Block) {
	api.AddBlockToPoolFunc(block)
}

type MockCommandHandlerPeerRepository struct {
}

func (rp MockCommandHandlerPeerRepository) Add(peer blockchain.Peer) error { return nil }
func (rp MockCommandHandlerPeerRepository) Remove(peer blockchain.PeerId) error { return nil }

type MockEventRepository struct {}

func (er MockEventRepository) Save(aggregateID string, events ...midgard.Event) error { return nil }
func (er MockEventRepository) Load(aggregate midgard.Aggregate, aggregateID string) error { return nil }
func (er MockEventRepository) Close() {}

func TestCommandHandler_HandleConfirmBlockCommand(t *testing.T) {
	tests := map[string] struct {
		input struct {
			command blockchain.ConfirmBlockCommand
		}
		err error
	} {
		"success": {
			input: struct {
				command blockchain.ConfirmBlockCommand
			} {
				command: blockchain.ConfirmBlockCommand{
					CommandModel: midgard.CommandModel{ID: "zf"},
					Block: &blockchain.DefaultBlock{
						Height: 99887,
					},
				},
			},
			err: nil,
		},
		"block nil error test": {
			input: struct {
				command blockchain.ConfirmBlockCommand
			} {
				command: blockchain.ConfirmBlockCommand{
					CommandModel: midgard.CommandModel{ID: "zf"},
					Block: nil,
				},
			},
			err: adapter.ErrBlockNil,
		},
	}

	blockApi := MockBlockApi{}
	blockApi.AddBlockToPoolFunc = func(block blockchain.Block) {
		assert.Equal(t, block.GetHeight(), uint64(99887))
	}

	peerRepository := MockPeerRepository{}
	eventRepository := MockEventRepository{}

	repositoryProjector := adapter.RepositoryProjector{
		PeerRepository: peerRepository,
		EventRepository: eventRepository,
	}

	commandHandler := adapter.NewCommandHandler(blockApi, repositoryProjector)
	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		err := commandHandler.HandleConfirmBlockCommand(test.input.command)

		assert.Equal(t, err, test.err)
	}

}
