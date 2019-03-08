/*
 * Copyright 2018 DE-labtory
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

	"github.com/DE-labtory/engine/api_gateway"
	"github.com/DE-labtory/engine/blockchain/api"
	"github.com/DE-labtory/engine/blockchain/infra/adapter"
	"github.com/DE-labtory/engine/blockchain/infra/mem"
	"github.com/DE-labtory/engine/blockchain/infra/repo"
	"github.com/DE-labtory/engine/common"
	"github.com/DE-labtory/engine/common/rabbitmq/pubsub"
	"github.com/DE-labtory/engine/conf"
	"github.com/DE-labtory/iLogger"
	"go.uber.org/fx"
)

const publisherID = "publisher.1"
const BbPath = "./db"

var Module = fx.Options(
	fx.Provide(
		NewBlockRepository,
		NewSyncStateRepository,
		mem.NewBlockPool,
		NewBlockAdapter,
		NewQueryService,
		NewBlockApi,
		NewSyncApi,
		NewConnectionEventHandler,
		NewBlockProposeHandler,
		NewConsensusEventHandler,
	),
	fx.Invoke(
		RegisterPubsubHandlers,
		RegisterTearDown,
		CreateGenesisBlock,
	),
)

func NewBlockAdapter() *adapter.HttpBlockAdapter {
	return adapter.NewHttpBlockAdapter()
}

func NewQueryService(blockAdapter *adapter.HttpBlockAdapter, peerQueryApi *api_gateway.PeerQueryApi) *adapter.QuerySerivce {
	return adapter.NewQueryService(blockAdapter, peerQueryApi)
}

func NewBlockRepository() (*repo.BlockRepository, error) {

	return repo.NewBlockRepository(BbPath)
}

func NewSyncStateRepository() *mem.SyncStateRepository {
	return mem.NewSyncStateRepository()
}

func NewBlockApi(config *conf.Configuration, blockRepository *repo.BlockRepository, blockPool *mem.BlockPool, service common.EventService, nodeId common.NodeID) (*api.BlockApi, error) {
	return api.NewBlockApi(nodeId, blockRepository, service, blockPool)
}

func NewSyncApi(config *conf.Configuration, blockRepository *repo.BlockRepository, syncStateRepository *mem.SyncStateRepository, eventService common.EventService, queryService *adapter.QuerySerivce, blockPool *mem.BlockPool, nodeId common.NodeID) (*api.SyncApi, error) {
	api, err := api.NewSyncApi(nodeId, blockRepository, syncStateRepository, eventService, queryService, blockPool)
	return &api, err
}

func NewBlockProposeHandler(blockApi *api.BlockApi, config *conf.Configuration) *adapter.BlockProposeCommandHandler {
	return adapter.NewBlockProposeCommandHandler(blockApi, config.Engine.Mode)
}

func NewConnectionEventHandler(syncApi *api.SyncApi) *adapter.NetworkEventHandler {
	return adapter.NewNetworkEventHandler(syncApi)
}

func NewConsensusEventHandler(syncStateRepository *mem.SyncStateRepository, blockApi *api.BlockApi) *adapter.ConsensusEventHandler {
	return adapter.NewConsensusEventHandler(syncStateRepository, blockApi)

}

func CreateGenesisBlock(blockApi *api.BlockApi, config *conf.Configuration) {
	if err := blockApi.CommitGenesisBlock(config.Blockchain.GenesisConfPath); err != nil {
		panic(err)
	}
}

func RegisterPubsubHandlers(subscriber *pubsub.TopicSubscriber, networkEventHandler *adapter.NetworkEventHandler, blockCommandHandler *adapter.BlockProposeCommandHandler, consensusEventHandler *adapter.ConsensusEventHandler) {
	iLogger.Infof(nil, "[Main] Blockchain is starting")

	if err := subscriber.SubscribeTopic("network.joined", networkEventHandler); err != nil {
		panic(err)
	}

	if err := subscriber.SubscribeTopic("block.propose", blockCommandHandler); err != nil {
		panic(err)
	}

	if err := subscriber.SubscribeTopic("block.confirm", consensusEventHandler); err != nil {
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
