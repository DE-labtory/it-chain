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

package api_gatewayfx

import (
	"context"
	"net/http"
	"os"

	"github.com/it-chain/iLogger"

	kitlog "github.com/go-kit/kit/log"
	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/conf"
	"go.uber.org/fx"
)

const ApidbPath = "./api-db"

var Module = fx.Options(
	fx.Provide(
		//api_gateway.NewConnectionRepository,
		NewKitLogger,
		NewBlockRepository,
		NewICodeRepository,
		NewPeerRepository,
		NewBlockQueryApi,
		NewBlockEventListener,
		api_gateway.NewConnectionEventListener,
		api_gateway.NewLeaderUpdateEventListener,
		NewICodeQueryApi,
		NewICodeEventHandler,
		api_gateway.NewPeerQueryApi,
		NewIvmHttpApi,
		api_gateway.NewApiHandler,
		http.NewServeMux,
	),
	fx.Invoke(
		RegisterEvent,
		RegisterHandlers,
		InitApiGatewayServer,
	),
)

func NewBlockRepository() (*api_gateway.BlockRepositoryImpl, error) {
	blockchainDB := ApidbPath + "/block"
	return api_gateway.NewBlockRepositoryImpl(blockchainDB)
}

func NewKitLogger() kitlog.Logger {
	var kitLogger kitlog.Logger
	kitLogger = kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stderr))
	kitLogger = kitlog.With(kitLogger, "ts", kitlog.DefaultTimestampUTC)
	httpLogger := kitlog.With(kitLogger, "component", "http")

	return httpLogger
}

func NewICodeQueryApi(repository *api_gateway.LevelDbICodeRepository) *api_gateway.ICodeQueryApi {
	return api_gateway.NewICodeQueryApi(repository)
}

func NewIvmHttpApi() *api_gateway.ICodeCommandApi {
	return api_gateway.NewICodeCommandApi()
}

func NewICodeEventHandler(repository *api_gateway.LevelDbICodeRepository) *api_gateway.ICodeEventHandler {
	return api_gateway.NewIcodeEventHandler(repository)
}

func NewBlockEventListener(blockRepository *api_gateway.BlockRepositoryImpl) *api_gateway.BlockEventListener {
	return api_gateway.NewBlockEventListener(blockRepository)
}

func NewBlockQueryApi(blockRepository *api_gateway.BlockRepositoryImpl) *api_gateway.BlockQueryApi {
	return api_gateway.NewBlockQueryApi(blockRepository)
}

func NewICodeRepository() *api_gateway.LevelDbICodeRepository {
	icodeDB := ApidbPath + "/ivm"
	return api_gateway.NewLevelDbMetaRepository(icodeDB)
}

func NewPeerRepository(config *conf.Configuration) *api_gateway.PeerRepository {

	var role api_gateway.Role
	if config.Engine.BootstrapNodeAddress == "" {
		role = api_gateway.Leader
	} else {
		role = api_gateway.Member
	}

	NodeId := common.GetNodeID(config.Engine.KeyPath, "ECDSA256")
	peerRepository := api_gateway.NewPeerRepository()
	peerRepository.Save(api_gateway.Peer{
		ID:                 NodeId,
		Role:               role,
		GrpcGatewayAddress: config.GrpcGateway.Address + ":" + config.GrpcGateway.Port,
		ApiGatewayAddress:  config.ApiGateway.Address + ":" + config.ApiGateway.Port,
	})

	return peerRepository
}

func RegisterEvent(subscriber *pubsub.TopicSubscriber, blockEventListener *api_gateway.BlockEventListener, icodeEventListener *api_gateway.ICodeEventHandler, connectionEventhandler *api_gateway.ConnectionEventHandler, leaderUpdateEventlistener *api_gateway.LeaderUpdateEventListener) {
	if err := subscriber.SubscribeTopic("block.*", blockEventListener); err != nil {
		panic(err)
	}
	if err := subscriber.SubscribeTopic("icode.*", icodeEventListener); err != nil {
		panic(err)
	}
	if err := subscriber.SubscribeTopic("connection.*", connectionEventhandler); err != nil {
		panic(err)
	}
	if err := subscriber.SubscribeTopic("leader.updated", leaderUpdateEventlistener); err != nil {
		panic(err)
	}
}

func RegisterHandlers(mux *http.ServeMux) {
	http.Handle("/", mux)
}

func InitApiGatewayServer(lifecycle fx.Lifecycle, config *conf.Configuration, handler http.Handler, blockRepo *api_gateway.BlockRepositoryImpl, iCodeRepo *api_gateway.LevelDbICodeRepository) {
	ipAddress := config.ApiGateway.Address + ":" + config.ApiGateway.Port

	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			iLogger.Infof(nil, "[Main] Api-gateway is staring on port:%s", config.ApiGateway.Port)
			go func() {
				err := http.ListenAndServe(ipAddress, handler)
				if err != nil {
					panic(err.Error())
				}
			}()
			return nil
		},
		OnStop: func(context context.Context) error {
			os.RemoveAll(ApidbPath)
			blockRepo.Close()
			iCodeRepo.Close()
			return nil
		},
	})
}
