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

package ivmfx

import (
	"context"

	"github.com/it-chain/iLogger"

	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/ivm"
	"github.com/it-chain/engine/ivm/api"
	"github.com/it-chain/engine/ivm/infra/adapter"
	"github.com/it-chain/engine/ivm/infra/git"
	"github.com/it-chain/engine/ivm/infra/tesseract"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewGitReposutoryService,
		NewContainerService,
		api.NewICodeApi,
		adapter.NewDeployCommandHandler,
		adapter.NewUnDeployCommandHandler,
		adapter.NewIcodeExecuteCommandHandler,
		adapter.NewListCommandHandler,
		adapter.NewBlockCommittedEventHandler,
	),
	fx.Invoke(
		RegisterRpcHandlers,
		RegisterPubsubHandlers,
		RegisterTearDown,
	),
)

func NewGitReposutoryService() ivm.GitService {
	return git.NewRepositoryService()
}

func NewContainerService() ivm.ContainerService {
	return tesseract.NewContainerService()
}

func RegisterRpcHandlers(
	server *rpc.Server,
	executeCommandHandler *adapter.IcodeExecuteCommandHandler,
	listCommandHandler *adapter.ListCommandHandler,
	deployCommandHandler *adapter.DeployCommandHandler,
	unDeployCommandHandler *adapter.UnDeployCommandHandler,
) {

	if err := server.Register("ivm.execute", executeCommandHandler.HandleTransactionExecuteCommandHandler); err != nil {
		panic(err)
	}
	if err := server.Register("ivm.deploy", deployCommandHandler.HandleDeployCommand); err != nil {
		panic(err)
	}
	if err := server.Register("ivm.undeploy", unDeployCommandHandler.HandleUnDeployCommand); err != nil {
		panic(err)
	}
	if err := server.Register("ivm.list", listCommandHandler.HandleListCommand); err != nil {
		panic(err)
	}
}

func RegisterPubsubHandlers(subscriber *pubsub.TopicSubscriber, handler *adapter.BlockCommittedEventHandler) {
	iLogger.Infof(nil, "[Main] Ivm is starting")
	if err := subscriber.SubscribeTopic("block.*", handler); err != nil {
		panic(err)
	}
}

func RegisterTearDown(lifecycle fx.Lifecycle, containerService ivm.ContainerService) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			return nil
		},
		OnStop: func(context context.Context) error {
			iCodeInfos := containerService.GetRunningICodeList()
			for _, iCodeInfo := range iCodeInfos {
				containerService.StopContainer(iCodeInfo.ID)
			}
			return nil
		},
	})
}
