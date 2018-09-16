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
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/api"
	"github.com/it-chain/engine/consensus/pbft/infra/adapter"
	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/infra/mem"
)

// 프로세스 아이디와 동일한 값의 ip를 가지는 프로세스들을 만들어낸다.
// set test environment that has ip address which is same as the process id
// 반환값은 모든 프로세스와 event service들이다.

func SetTestEnvironment(processList []string) struct {
	ProcessMap      map[string]*mock.Process
	EventServiceMap map[string]*mock.EventService
} {
	networkManager := mock.NewNetworkManager()

	processMap := make(map[string]*mock.Process)
	eventServiceMap := make(map[string]*mock.EventService)

	for _, id := range processList {

		// setup process
		process := mock.NewProcess(id)
		electionService := pbft.NewElectionService(id, 30, pbft.TICKING, 0)
		repository := mem.NewPeerReopository()
		for _, pid := range processList {
			repository.Save(p2p.Peer{
				PeerId: p2p.PeerId{
					Id: pid,
				},
				IpAddress: pid,
			})
		}

		peerQueryApi := api_gateway.NewPeerQueryApi(repository)
		parliament := pbft.NewParliament()
		parliamentService := adapter.NewParliamentService(parliament, peerQueryApi)
		parliamentService.Build()
		eventService := mock.NewEventService(id, networkManager.Publish)
		electionApi := api.NewElectionApi(electionService, parliamentService, eventService)
		leaderApi := api.NewLeaderApi(parliamentService, eventService)
		grpcCommandHandler := adapter.NewElectionCommandHandler(leaderApi, electionApi)

		// register handler to process
		process.RegisterHandler(grpcCommandHandler.HandleMessageReceive)

		// register module to process
		process.Register(electionApi)
		process.Register(leaderApi)
		process.Register(electionService)

		// add process to network manager
		networkManager.AddProcess(process)
		processMap[process.Id] = process
		eventServiceMap[process.Id] = eventService
	}

	logger.Infof(nil, "created process: %v", processMap)
	logger.Infof(nil, "created event service: %v", eventServiceMap)

	networkManager.Start()
	return struct {
		ProcessMap      map[string]*mock.Process
		EventServiceMap map[string]*mock.EventService
	}{
		ProcessMap:      processMap,
		EventServiceMap: eventServiceMap,
	}
}
