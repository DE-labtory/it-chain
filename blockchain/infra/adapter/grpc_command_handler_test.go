package adapter_test

import (
	"errors"
	"testing"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/infra/adapter"
	"github.com/stretchr/testify/assert"
	"errors"
)

type MockSyncBlockApi struct {}
func (ba MockSyncBlockApi) SyncedCheck(block blockchain.Block) error { return nil }

type MockBlockQueryApi struct {
	GetLastBlockFunc func() (blockchain.Block, error)
}
func (br MockBlockQueryApi) GetLastBlock() (blockchain.Block, error) {
	return br.GetLastBlockFunc()
}
func (br MockBlockQueryApi) GetBlockByHeight(blockHeight uint64) (blockchain.Block, error) { return nil, nil }
func (br MockBlockQueryApi) GetBlockBySeal(seal []byte) (blockchain.Block, error) { return nil, nil }
func (br MockBlockQueryApi) GetBlockByTxID(txid string) (blockchain.Block, error) { return nil, nil }
func (br MockBlockQueryApi) GetTransactionByTxID(txid string) (blockchain.Transaction, error) { return nil, nil }

func (br MockROBlockRepository) GetBlockByHeight(block blockchain.Block, height uint64) error {
	return br.GetBlockByHeightFunc(block, height)
}

type MockGrpcCommandService struct {
	SyncCheckResponseFunc func(block blockchain.Block) error
	ResponseBlockFunc     func(peerId blockchain.PeerId, block blockchain.Block) error
}

func (cs MockGrpcCommandService) SyncCheckResponse(block blockchain.Block) error {
	return cs.SyncCheckResponseFunc(block)
}

func (cs MockGrpcCommandService) ResponseBlock(peerId blockchain.PeerId, block blockchain.Block) error {
	return cs.ResponseBlockFunc(peerId, block)
}

func TestGrpcCommandHandler_HandleGrpcCommand_SyncCheckRequestProtocol(t *testing.T) {
	tests := map[string]struct {
		input struct {
			command         blockchain.GrpcReceiveCommand
			getLastBlockErr error
			syncCheckErr    error
		}
		err error
	}{
		"success": {
			input: struct {
				command         blockchain.GrpcReceiveCommand
				getLastBlockErr error
				syncCheckErr    error
			}{
				command: blockchain.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{ID: "111"},
					Body:         nil,
					Protocol:     "SyncCheckRequestProtocol",
				},
			},
			err: nil,
		},
		"get last block err test": {
			input: struct {
				command         blockchain.GrpcReceiveCommand
				getLastBlockErr error
				syncCheckErr    error
			}{
				command: blockchain.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{ID: "111"},
					Body:         nil,
					Protocol:     "SyncCheckRequestProtocol",
				},
				getLastBlockErr: errors.New("error occur in ErrGetLastBlock"),
			},
			err: adapter.ErrGetLastBlock,
		},
		"sync check err test": {
			input: struct {
				command         blockchain.GrpcReceiveCommand
				getLastBlockErr error
				syncCheckErr    error
			}{
				command: blockchain.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{ID: "111"},
					Body:         nil,
					Protocol:     "SyncCheckRequestProtocol",
				},
				syncCheckErr: errors.New("error occur in SyncCheckResponse"),
			},
			err: adapter.ErrSyncCheckResponse,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		blockApi := MockSyncBlockApi{}

		blockRepository := MockBlockQueryApi{}
		blockRepository.GetLastBlockFunc = func() (blockchain.Block, error) {
			return &blockchain.DefaultBlock{Height: blockchain.BlockHeight(99887)},test.input.getLastBlockErr
		}

		grpcCommandService := MockGrpcCommandService{}
		grpcCommandService.SyncCheckResponseFunc = func(block blockchain.Block) error {
			assert.Equal(t, block.GetHeight(), uint64(99887))
			return test.input.syncCheckErr
		}

		grpcCommandHandler := adapter.NewGrpcCommandHandler(blockApi, blockRepository, grpcCommandService)

		err := grpcCommandHandler.HandleGrpcCommand(test.input.command)
		assert.Equal(t, err, test.err)
	}
}

func TestGrpcCommandHandler_HandleGrpcCommand_BlockRequestProtocol(t *testing.T) {
	tests := map[string]struct {
		input struct {
			command blockchain.GrpcReceiveCommand
			err     struct {
				ErrGetBlock      error
				ErrResponseBlock error
			}
		}
		err error
	}{
		"success:": {
			input: struct {
				command blockchain.GrpcReceiveCommand
				err     struct {
					ErrGetBlock      error
					ErrResponseBlock error
				}
			}{
				command: blockchain.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{ID: "111"},
					Body:         []byte{48},
					Protocol:     "BlockRequestProtocol",
				},

				err: struct {
					ErrGetBlock      error
					ErrResponseBlock error
				}{
					ErrGetBlock:      nil,
					ErrResponseBlock: nil,
				},
			},
			err: nil,
		},
		"fail: Umnarshal command": {
			input: struct {
				command blockchain.GrpcReceiveCommand
				err     struct {
					ErrGetBlock      error
					ErrResponseBlock error
				}
			}{
				command: blockchain.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{ID: "111"},
					Body:         nil,
					Protocol:     "BlockRequestProtocol",
				},

				err: struct {
					ErrGetBlock      error
					ErrResponseBlock error
				}{
					ErrGetBlock:      nil,
					ErrResponseBlock: nil,
				},
			},
			err: adapter.ErrBlockInfoDeliver,
		},

		"fail: get block by height error test": {
			input: struct {
				command blockchain.GrpcReceiveCommand
				err     struct {
					ErrGetBlock      error
					ErrResponseBlock error
				}
			}{
				command: blockchain.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{ID: "111"},
					Body:         []byte{48},
					Protocol:     "BlockRequestProtocol",
				},

				err: struct {
					ErrGetBlock      error
					ErrResponseBlock error
				}{
					ErrGetBlock:      errors.New("error when getting block by height"),
					ErrResponseBlock: nil,
				},
			},
			err: adapter.ErrGetBlock,
		},

		"fail: response block error test": {
			input: struct {
				command blockchain.GrpcReceiveCommand
				err     struct {
					ErrGetBlock      error
					ErrResponseBlock error
				}
			}{
				command: blockchain.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{ID: "111"},
					Body:         []byte{48},
					Protocol:     "BlockRequestProtocol",
				},

				err: struct {
					ErrGetBlock      error
					ErrResponseBlock error
				}{
					ErrGetBlock:      nil,
					ErrResponseBlock: errors.New("error when response block"),
				},
			},
			err: adapter.ErrResponseBlock,
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		blockApi := MockBlockApi{}

		blockRepository := MockROBlockRepository{}
		blockRepository.GetBlockByHeightFunc = func(block blockchain.Block, height uint64) error {
			block.SetHeight(height)
			return test.input.err.ErrGetBlock
		}

		grpcCommandService := MockGrpcCommandService{}
		grpcCommandService.ResponseBlockFunc = func(peerId blockchain.PeerId, block blockchain.Block) error {
			return test.input.err.ErrResponseBlock
		}

		grpcCommandHandler := adapter.NewGrpcCommandHandler(blockApi, blockRepository, grpcCommandService)

		err := grpcCommandHandler.HandleGrpcCommand(test.input.command)
		assert.Equal(t, err, test.err)

	}

}
