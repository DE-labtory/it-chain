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

	kitlog "github.com/go-kit/kit/log"
	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/conf"
	"go.uber.org/fx"
)

const ApidbPath = "./api-db"

var Module = fx.Options(
	fx.Provide(
		NewKitLogger,
		NewBlockRepository,
		NewICodeRepository,
		api_gateway.NewBlockQueryApi,
		api_gateway.NewBlockEventListener,
		api_gateway.NewICodeQueryApi,
		api_gateway.NewIcodeEventHandler,
		api_gateway.NewApiHandler,
		http.NewServeMux,
	),
	fx.Invoke(
		RegisterEvent,
		RegisterHandlers,
		InitApiGatewayServer,
	),
)

func NewBlockRepository() (api_gateway.BlockRepository, error) {
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

func NewICodeRepository() api_gateway.ICodeRepository {
	icodeDB := ApidbPath + "/ivm"
	return api_gateway.NewLevelDbMetaRepository(icodeDB)
}

func RegisterEvent(subscriber *pubsub.TopicSubscriber, blockEventListener *api_gateway.BlockEventListener, icodeEventListener *api_gateway.ICodeEventHandler) {
	if err := subscriber.SubscribeTopic("block.*", blockEventListener); err != nil {
		panic(err)
	}
	if err := subscriber.SubscribeTopic("icode.*", icodeEventListener); err != nil {
		panic(err)
	}
}

func RegisterHandlers(mux *http.ServeMux) {
	http.Handle("/", mux)
}

func InitApiGatewayServer(lifecycle fx.Lifecycle, config *conf.Configuration, handler http.Handler, blockRepo api_gateway.BlockRepository, iCodeRepo api_gateway.ICodeRepository) {
	ipAddress := config.ApiGateway.Address + ":" + config.ApiGateway.Port

	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			logger.Infof(nil, "[Main] Api-gateway is staring on port:%s", config.ApiGateway.Port)
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
