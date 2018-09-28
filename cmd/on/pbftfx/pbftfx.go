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

package pbftfx

import (
	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/api"
	"github.com/it-chain/engine/consensus/pbft/infra/adapter"
	"github.com/it-chain/iLogger"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		pbft.NewParliament,
		NewElectionService,
		NewParliamentService,
		NewElectionApi,
		NewLeaderApi,
		adapter.NewElectionCommandHandler,
		adapter.NewConnectionCreatedEventHandler,
	),
	fx.Invoke(
		RegisterPubsubHandlers,
	),
)

func NewElectionService(config *conf.Configuration) *pbft.ElectionService {
	NodeId := common.GetNodeID(config.Engine.KeyPath, "ECDSA256")
	return pbft.NewElectionService(NodeId, 30, pbft.CANDIDATE, 0)
}

func NewParliamentService(parliament *pbft.Parliament, peerQueryApi *api_gateway.PeerQueryApi) *adapter.ParliamentService {
	return adapter.NewParliamentService(parliament, peerQueryApi)
}

func NewElectionApi(electionService *pbft.ElectionService, parliamentService *adapter.ParliamentService, eventService common.EventService) *api.ElectionApi {
	return api.NewElectionApi(electionService, parliamentService, eventService)
}

func NewLeaderApi(parliamentService *adapter.ParliamentService, eventService common.EventService) *api.LeaderApi {
	return api.NewLeaderApi(parliamentService, eventService)
}

func RegisterPubsubHandlers(subscriber *pubsub.TopicSubscriber, electionCommandHandler *adapter.ElectionCommandHandler, connectionCreatedEventHandler *adapter.ConnectionCreatedEventHandler) {
	iLogger.Infof(nil, "[Main] Consensus is starting")

	if err := subscriber.SubscribeTopic("message.*", electionCommandHandler); err != nil {
		panic(err)
	}

	if err := subscriber.SubscribeTopic("connection.created", connectionCreatedEventHandler); err != nil {
		panic(err)
	}
}
