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
	"errors"
	"testing"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/it-chain/midgard"
	"github.com/stretchr/testify/assert"
)

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

		blockApi := mock.MockSyncBlockApi{}

		blockRepository := mock.BlockQueryApi{}
		blockRepository.GetLastBlockFunc = func() (blockchain.Block, error) {
			return &blockchain.DefaultBlock{Height: blockchain.BlockHeight(99887)}, test.input.getLastBlockErr
		}

		grpcCommandService := mock.SyncCheckGrpcCommandService{}
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

		blockApi := mock.MockSyncBlockApi{}

		blockQueryApi := mock.BlockQueryApi{}
		blockQueryApi.GetBlockByHeightFunc = func(height uint64) (blockchain.Block, error) {
			return &blockchain.DefaultBlock{
				Height: blockchain.BlockHeight(12),
			}, test.input.err.ErrGetBlock
		}

		grpcCommandService := mock.SyncCheckGrpcCommandService{}
		grpcCommandService.ResponseBlockFunc = func(peerId blockchain.PeerId, block blockchain.Block) error {
			return test.input.err.ErrResponseBlock
		}

		grpcCommandHandler := adapter.NewGrpcCommandHandler(blockApi, blockQueryApi, grpcCommandService)

		err := grpcCommandHandler.HandleGrpcCommand(test.input.command)
		assert.Equal(t, err, test.err)

	}

}
