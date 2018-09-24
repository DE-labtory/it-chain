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

package grpc_gatewayfx

import (
	"context"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/grpc_gateway/api"
	"github.com/it-chain/engine/grpc_gateway/infra"
	"github.com/it-chain/engine/grpc_gateway/infra/adapter"
	"github.com/it-chain/iLogger"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		adapter.NewHttpPeerAdapter,
		NewPeerService,
		NewGrpcHostService,
		NewConnectionApi,
		adapter.NewConnectionCommandHandler,
		adapter.NewGrpcMessageHandler,
	),
	fx.Invoke(
		RegisterHandlers,
		RegisterEvent,
		InitgRPCServer,
	),
)

func NewPeerService(peerAdapter *adapter.HttpPeerAdapter) *adapter.PeerService {
	return adapter.NewPeerService(peerAdapter)
}

func NewGrpcHostService(conf *conf.Configuration, publisher *pubsub.TopicPublisher) *infra.GrpcHostService {
	priKey, pubKey := infra.LoadKeyPair(conf.Engine.KeyPath, "ECDSA256")
	hostService := infra.NewGrpcHostService(priKey, pubKey, publisher.Publish, infra.HostInfo{
		ApiGatewayAddress:  conf.ApiGateway.Address + ":" + conf.ApiGateway.Port,
		GrpcGatewayAddress: conf.GrpcGateway.Address + ":" + conf.GrpcGateway.Port,
	})
	return hostService
}

func NewConnectionApi(hostService *infra.GrpcHostService, eventService common.EventService, peerService *adapter.PeerService) *api.ConnectionApi {
	return api.NewConnectionApi(hostService, eventService, peerService)
}

func RegisterHandlers(connectionCommandHandler *adapter.ConnectionCommandHandler, server *rpc.Server) {
	iLogger.Infof(nil, "[Main] gRPC-Gateway is starting")
	if err := server.Register("connection.create", connectionCommandHandler.HandleCreateConnectionCommand); err != nil {
		panic(err)
	}

	if err := server.Register("connection.list", connectionCommandHandler.HandleGetConnectionListCommand); err != nil {
		panic(err)
	}

	if err := server.Register("connection.close", connectionCommandHandler.HandleCloseConnectionCommand); err != nil {
		panic(err)
	}

	if err := server.Register("connection.join", connectionCommandHandler.HandleJoinNetworkCommand); err != nil {
		panic(err)
	}
}

func RegisterEvent(grpcCommandHandler *adapter.GrpcMessageHandler, subscriber *pubsub.TopicSubscriber) {
	if err := subscriber.SubscribeTopic("message.receive", grpcCommandHandler); err != nil {
		panic(err)
	}
}

func InitgRPCServer(lifecycle fx.Lifecycle, config *conf.Configuration, hostService *infra.GrpcHostService, connectionApi *api.ConnectionApi) {
	hostService.SetHandler(connectionApi)

	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			go hostService.Listen(config.GrpcGateway.Address + ":" + config.GrpcGateway.Port)
			return nil
		},
		OnStop: func(context context.Context) error {
			hostService.CloseAllConnections()
			hostService.Stop()
			return nil
		},
	})
}
