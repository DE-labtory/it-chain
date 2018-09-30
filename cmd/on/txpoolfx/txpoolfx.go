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

package txpoolfx

import (
	"context"
	"time"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/txpool/api"
	"github.com/it-chain/engine/txpool/infra/adapter"
	"github.com/it-chain/engine/txpool/infra/batch"
	"github.com/it-chain/engine/txpool/infra/mem"
	"github.com/it-chain/iLogger"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		mem.NewTransactionRepository,
		NewBlockProposalService,
		NewTxpoolApi,
		adapter.NewTxCommandHandler,
	),
	fx.Invoke(
		RunBatcher,
		RegisterRpcHandlers,
	),
)

func NewBlockProposalService(repository *mem.TransactionRepository, client *rpc.Client, config *conf.Configuration) *adapter.BlockProposalService {
	return adapter.NewBlockProposalService(client, repository, config.Engine.Mode)
}

func NewTxpoolApi(config *conf.Configuration, repository *mem.TransactionRepository) *api.TransactionApi {
	NodeId := common.GetNodeID(config.Engine.KeyPath, "ECDSA256")
	return api.NewTransactionApi(NodeId, repository)
}

func RunBatcher(lifecycle fx.Lifecycle, blockProposalService *adapter.BlockProposalService, config *conf.Configuration) {
	var q chan struct{}
	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			q = batch.GetTimeOutBatcherInstance().Run(blockProposalService.ProposeBlock, (time.Duration(config.Txpool.TimeoutMs) * time.Millisecond))
			return nil
		},
		OnStop: func(context context.Context) error {
			q <- struct{}{}
			return nil
		},
	})
}

func RegisterRpcHandlers(server *rpc.Server, handler *adapter.TxCommandHandler) {
	iLogger.Infof(nil, "[Main] Txpool is starting")
	if err := server.Register("transaction.create", handler.HandleTxCreateCommand); err != nil {
		panic(err)
	}
}
