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

package test

import (
	"github.com/it-chain/avengers/mock"
	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/api"
	"github.com/it-chain/engine/p2p/infra/adapter"
	"github.com/it-chain/engine/p2p/infra/mem"
	mock2 "github.com/it-chain/engine/p2p/test/mock"
)

func SetTestEnvironment(processList []string) map[string]*mock.Process {
	networkManager := mock.NewNetworkManager()

	m := make(map[string]*mock.Process)
	for _, processId := range processList {
		process := mock.NewProcess()
		process.Init(processId)
		networkManager.AddProcess(process)

		election := p2p.NewElection(processId, 30, p2p.Ticking, 0)
		peerRepository := mem.NewPeerReopository()

		peerQueryService := api_gateway.NewPeerQueryApi(&peerRepository)
		client := mock.NewClient(processId, networkManager.GrpcCall)
		server := mock.NewServer(processId, networkManager.GrpcConsume)

		eventService := mock2.MockEventService{}

		eventService.PublishFunc = func(topic string, event interface{}) error {
			return nil
		}

		electionService := p2p.NewElectionService(&election, &peerQueryService, &client)

		pLTableService := p2p.PLTableService{}

		communicationService := p2p.NewCommunicationService(&client)

		communicationApi := api.NewCommunicationApi(&peerQueryService, communicationService)

		leaderApi := api.NewLeaderApi(&peerRepository, &eventService)
		grpcCommandHandler := adapter.NewGrpcCommandHandler(&leaderApi, &electionService, &communicationApi, pLTableService)
		server.Register("message.receive", grpcCommandHandler.HandleMessageReceive)

		process.Register(&electionService)
		process.Register(&peerRepository)
		m[process.Id] = &process
	}

	logger.Infof(nil, "created process: %v", m)

	return m
}
