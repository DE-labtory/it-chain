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
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/rabbitmq/pubsub"
	"github.com/it-chain/engine/conf"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/api"
	"github.com/it-chain/engine/consensus/pbft/infra/adapter"
	"github.com/it-chain/engine/consensus/pbft/infra/mem"
	"github.com/it-chain/iLogger"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		NewParliamentRepository,
		NewElectionService,
		NewElectionApi,
		NewParliamentApi,
		adapter.NewElectionCommandHandler,
		adapter.NewConnectionEventHandler,
		adapter.NewLeaderCommandHandler,
	),
	fx.Invoke(
		RegisterPubsubHandlers,
	),
)

func NewElectionService(config *conf.Configuration) *pbft.ElectionService {
	NodeId := common.GetNodeID(config.Engine.KeyPath, "ECDSA256")
	return pbft.NewElectionService(NodeId, 30, pbft.TICKING, 0)
}

func NewParliamentRepository(config *conf.Configuration) *mem.ParliamentRepository {

	NodeId := common.GetNodeID(config.Engine.KeyPath, "ECDSA256")
	parliament := pbft.NewParliament()
	parliament.AddRepresentative(pbft.NewRepresentative(NodeId))

	if config.Engine.BootstrapNodeAddress == "" {
		parliament.SetLeader(NodeId)
	}

	return mem.NewParliamentRepositoryWithParliament(parliament)
}

func NewElectionApi(electionService *pbft.ElectionService, parliamentRepository *mem.ParliamentRepository, eventService common.EventService) *api.ElectionApi {
	return api.NewElectionApi(electionService, parliamentRepository, eventService)
}

func NewParliamentApi(config *conf.Configuration, parliamentRepository *mem.ParliamentRepository, eventService common.EventService) *api.ParliamentApi {
	NodeId := common.GetNodeID(config.Engine.KeyPath, "ECDSA256")

	return api.NewParliamentApi(NodeId, parliamentRepository, eventService)
}

func RegisterPubsubHandlers(subscriber *pubsub.TopicSubscriber, electionCommandHandler *adapter.ElectionCommandHandler, connectionEventHandler *adapter.ConnectionEventHandler, leaderCommandHandler *adapter.LeaderCommandHandler) {
	iLogger.Infof(nil, "[Main] Consensus is starting")

	if err := subscriber.SubscribeTopic("message.receive", electionCommandHandler); err != nil {
		panic(err)
	}

	if err := subscriber.SubscribeTopic("connection.created", connectionEventHandler); err != nil {
		panic(err)
	}

	if err := subscriber.SubscribeTopic("message.receive", leaderCommandHandler); err != nil {
		panic(err)
	}
}
