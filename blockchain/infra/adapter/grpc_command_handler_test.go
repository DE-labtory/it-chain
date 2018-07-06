package adapter_test

import (
	"errors"
	"testing"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/infra/adapter"
	"github.com/it-chain/midgard"
	"github.com/magiconair/properties/assert"
)

type MockBlockApi struct{}

func (ba MockBlockApi) SyncedCheck(block blockchain.Block) error { return nil }

type MockROBlockRepository struct {
	NewEmptyBlockFunc    func() (blockchain.Block, error)
	GetLastBlockFunc     func(block blockchain.Block) error
	GetBlockByHeightFunc func(block blockchain.Block, height uint64) error
}

func (br MockROBlockRepository) NewEmptyBlock() (blockchain.Block, error) {
	return br.NewEmptyBlockFunc()
}
func (br MockROBlockRepository) GetLastBlock(block blockchain.Block) error {
	return br.GetLastBlockFunc(block)
}

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

		blockApi := MockBlockApi{}

		blockRepository := MockROBlockRepository{}
		blockRepository.GetLastBlockFunc = func(block blockchain.Block) error {
			block.SetHeight(99887)
			return test.input.getLastBlockErr
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

//func TestGrpcCommandHandler_HandleGrpcCommand_BlockRequestProtocol(t *testing.T) {
//	tests := map[string]struct {
//		input struct {
//			command blockchain.GrpcReceiveCommand
//		}
//		err error
//	}{
//		"success:": {
//			input: struct {
//				command blockchain.GrpcReceiveCommand
//			}{
//				command: blockchain.GrpcReceiveCommand{
//					CommandModel: midgard.CommandModel{ID: "111"},
//					Body:         nil,
//					Protocol:     "RequestBlockProtocol",
//				},
//			},
//			err: nil,
//		},
//
//		"fail:": {
//			input: struct {
//				command blockchain.GrpcReceiveCommand
//			}{
//				command: blockchain.GrpcReceiveCommand{
//					CommandModel: midgard.CommandModel{ID: "111"},
//					Body:         nil,
//					Protocol:     "RequestBlockProtocol",
//				},
//			},
//			err: nil,
//		},
//	}
//
//	for testName, test := tests {
//		t.Logf("running test case %s", testName)
//
//		blockApi := MockBlockApi{}
//
//		blockRepository := MockROBlockRepository{}
//		blockRepository.GetBlockByHeightFunc = func(block blockchain.Block,height uint64) error {
//
//		}
//
//
//
//	}
//
//}
