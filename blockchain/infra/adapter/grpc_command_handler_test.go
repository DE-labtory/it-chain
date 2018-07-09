package adapter_test

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"testing"
	"github.com/it-chain/midgard"
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

type MockSyncCheckGrpcCommandService struct {
	SyncCheckResponseFunc func(block blockchain.Block) error
}
func (cs MockSyncCheckGrpcCommandService) SyncCheckResponse(block blockchain.Block) error {
	return cs.SyncCheckResponseFunc(block)
}

func TestGrpcCommandHandler_HandleGrpcCommand_SyncCheckRequestProtocol(t *testing.T) {
	tests := map[string]struct {
		input struct {
			command blockchain.GrpcReceiveCommand
			getLastBlockErr error
			syncCheckErr error
		}
		err error
	} {
		"success": {
			input: struct {
				command blockchain.GrpcReceiveCommand
				getLastBlockErr error
				syncCheckErr error
			} {
				command: blockchain.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{ID: "111"},
					Body: nil,
					Protocol: "SyncCheckRequestProtocol",
				},
			},
			err: nil,
		},
		"get last block err test": {
			input: struct {
				command blockchain.GrpcReceiveCommand
				getLastBlockErr error
				syncCheckErr error
			} {
				command: blockchain.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{ID: "111"},
					Body: nil,
					Protocol: "SyncCheckRequestProtocol",
				},
				getLastBlockErr: errors.New("error occur in ErrGetLastBlock"),
			},
			err: adapter.ErrGetLastBlock,
		},
		"sync check err test": {
			input: struct {
				command blockchain.GrpcReceiveCommand
				getLastBlockErr error
				syncCheckErr error
			} {
				command: blockchain.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{ID: "111"},
					Body: nil,
					Protocol: "SyncCheckRequestProtocol",
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

		grpcCommandService := MockSyncCheckGrpcCommandService{}
		grpcCommandService.SyncCheckResponseFunc = func(block blockchain.Block) error {
			assert.Equal(t, block.GetHeight(), uint64(99887))
			return test.input.syncCheckErr
		}

		grpcCommandHandler := adapter.NewGrpcCommandHandler(blockApi, blockRepository, grpcCommandService)


		err := grpcCommandHandler.HandleGrpcCommand(test.input.command)
		assert.Equal(t, err, test.err)
	}
}
