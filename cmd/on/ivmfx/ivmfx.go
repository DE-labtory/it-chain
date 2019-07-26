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

package ivmfx

import (
	"context"

	"github.com/DE-labtory/iLogger"
	"github.com/DE-labtory/it-chain/common/rabbitmq/pubsub"
	"github.com/DE-labtory/it-chain/common/rabbitmq/rpc"
	"github.com/DE-labtory/it-chain/conf"
	"github.com/DE-labtory/it-chain/ivm"
	"github.com/DE-labtory/it-chain/ivm/api"
	"github.com/DE-labtory/it-chain/ivm/infra/adapter"
	"github.com/DE-labtory/it-chain/ivm/infra/git"
	"github.com/DE-labtory/it-chain/ivm/infra/tesseract"
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

func NewContainerService(conf *conf.Configuration) ivm.ContainerService {
	if conf.Docker.Use {
		cs, err := tesseract.NewContainerService(&tesseract.ContainerDockerConfig{
			Subnet:      conf.Docker.NetworkSubnet,
			VolumeName:  conf.Docker.VolumeName,
			NetworkName: conf.Docker.NetworkName,
		})
		if err != nil {
			panic(err)
		}
		return cs
	} else {
		cs, err := tesseract.NewContainerService(nil)
		if err != nil {
			panic(err)
		}
		return cs
	}

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
