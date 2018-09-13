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

package blockchainfx

import (
	"context"
	"os"

	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/engine/blockchain/api"
	"github.com/it-chain/engine/blockchain/infra/adapter"
	"github.com/it-chain/engine/blockchain/infra/mem"
	"github.com/it-chain/engine/blockchain/test/mock"
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"go.uber.org/fx"
)

const publisherID = "publisher.1"
const BbPath = "./db"

var Module = fx.Options(
	fx.Provide(
		NewBlockRepository,
		blockchain.NewBlockPool,
		NewBlockApi,
		NewBlockProposeHandler,
	),
	fx.Invoke(
		CreateGenesisBlock,
		RegisterRpcHandlers,
		RegisterTearDown,
	),
)

func NewBlockRepository() (*mem.BlockRepository, error) {
	return mem.NewBlockRepository(BbPath)
}

func NewBlockApi(blockRepository *mem.BlockRepository, blockPool *blockchain.BlockPool, service common.EventService) (*api.BlockApi, error) {
	return api.NewBlockApi(publisherID, blockRepository, service, blockPool)
}

func NewBlockProposeHandler(blockApi *api.BlockApi, config *conf.Configuration) *adapter.BlockProposeCommandHandler {
	return adapter.NewBlockProposeCommandHandler(blockApi, mock.ConsensusService{}, config.Engine.Mode)
}

func CreateGenesisBlock(blockApi *api.BlockApi, config *conf.Configuration) {
	if err := blockApi.CommitGenesisBlock(config.Blockchain.GenesisConfPath); err != nil {
		panic(err)
	}
}

func RegisterRpcHandlers(server *rpc.Server, handler *adapter.BlockProposeCommandHandler) {
	logger.Infof(nil, "[Main] Blockchain is starting")
	if err := server.Register("block.propose", handler.HandleProposeBlockCommand); err != nil {
		panic(err)
	}
}

func RegisterTearDown(lifecycle fx.Lifecycle) {

	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			return nil
		},
		OnStop: func(context context.Context) error {
			return os.RemoveAll(BbPath)
		},
	})
}
