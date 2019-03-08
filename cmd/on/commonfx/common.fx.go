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

package commonfx

import (
	"context"

	"github.com/DE-labtory/engine/common"
	"github.com/DE-labtory/engine/common/rabbitmq/pubsub"
	"github.com/DE-labtory/engine/common/rabbitmq/rpc"
	"github.com/DE-labtory/engine/conf"
	"github.com/DE-labtory/heimdall"
	"github.com/DE-labtory/heimdall/config"
	"github.com/DE-labtory/iLogger"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewConfiguration,
		NewSecConfig,
		NewKeyPair,
		NewNodeID,
		NewRpcServer,
		NewRpcClient,
		NewPubsubServer,
		NewEventService,
		NewPubsubClient,
	),
	fx.Invoke(
		RegisterTearDown,
	),
)

func NewConfiguration() *conf.Configuration {
	return conf.GetConfiguration()
}

func NewSecConfig(conf *conf.Configuration) *config.Config {
	secConf, err := config.NewSimpleConfig(conf.Engine.SecLv)
	if err != nil {
		panic(err)
	}
	return secConf
}

//[INFO] This function is temporary, will be removed when mock certs and keys added with cert related functions
//NewKeyPair generates key pair by configuration
func NewKeyPair(conf *conf.Configuration, secConf *config.Config) (*heimdall.PriKey, *heimdall.PubKey) {
	pri, pub, err := common.GenerateAndStoreKeyPair(secConf, conf.Engine.PriKeyPath)
	if err != nil {
		iLogger.Errorf(nil, "%s", err)
		panic(err)
	}

	return &pri, &pub
}

func NewNodeID(config *conf.Configuration, secConf *config.Config, pubKey *heimdall.PubKey) common.NodeID {
	return common.GetNodeID(secConf.SigAlgo, config.Engine.PriKeyPath)
}

func NewRpcServer(config *conf.Configuration) *rpc.Server {
	return rpc.NewServer(config.Engine.Amqp)
}

func NewRpcClient(config *conf.Configuration) *rpc.Client {
	return rpc.NewClient(config.Engine.Amqp)
}

func NewPubsubServer(config *conf.Configuration) *pubsub.TopicSubscriber {
	return pubsub.NewTopicSubscriber(config.Engine.Amqp, "Event")
}

func NewPubsubClient(config *conf.Configuration) *pubsub.TopicPublisher {
	return pubsub.NewTopicPublisher(config.Engine.Amqp, "Event")
}

func NewEventService(config *conf.Configuration) common.EventService {
	return common.NewEventService(config.Engine.Amqp, "Event")
}

func RegisterTearDown(lifecycle fx.Lifecycle, rpcServer *rpc.Server, subscriber *pubsub.TopicSubscriber, eventService common.EventService) {
	lifecycle.Append(fx.Hook{
		OnStart: func(context context.Context) error {
			return nil
		},
		OnStop: func(context context.Context) error {
			subscriber.Close()
			rpcServer.Close()
			eventService.Close()
			return nil
		},
	})
}
