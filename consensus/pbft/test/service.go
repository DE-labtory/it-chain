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
	"github.com/it-chain/engine/common/mock"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/api"
	"github.com/it-chain/engine/consensus/pbft/infra/adapter"
	"github.com/it-chain/engine/consensus/pbft/infra/mem"
	"github.com/it-chain/sdk/logger"
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
		electionService := pbft.NewElectionService(id, 30, pbft.NORMAL, 0)
		parliamentRepository := mem.NewParliamentRepository()
		parliament := pbft.NewParliament()

		for _, pid := range processList {
			parliament.AddRepresentative(pbft.NewRepresentative(pid))
		}

		parliamentRepository.Save(parliament)
		stateRepository := mem.NewStateRepository()

		eventService := mock.NewEventService(id, networkManager.Publish)
		propagateService := pbft.NewPropagateService(eventService)

		electionApi := api.NewElectionApi(electionService, parliamentRepository, eventService)
		leaderApi := api.NewParliamentApi(id, parliamentRepository, eventService)

		stateApi := api.NewStateApi(id, propagateService, eventService, parliamentRepository, stateRepository)

		grpcCommandHandler := adapter.NewElectionCommandHandler(leaderApi, electionApi)
		pbftHandler := adapter.NewPbftMsgHandler(stateApi)

		// register handler to process
		process.RegisterHandler(grpcCommandHandler.HandleMessageReceive)
		process.RegisterHandler(pbftHandler.HandleGrpcMsgCommand)

		// register module to process
		process.Register(electionApi)
		process.Register(leaderApi)
		process.Register(electionService)
		process.Register(stateApi)
		process.Register(parliamentRepository)
		process.Register(stateRepository)

		// add process to network manager
		networkManager.AddProcess(process)
		processMap[process.Id] = process
		eventServiceMap[process.Id] = eventService
	}

	logger.Infof(nil, "[consensus] created process: %v", processMap)
	logger.Infof(nil, "[consensus] created event service: %v", eventServiceMap)

	networkManager.Start()
	return struct {
		ProcessMap      map[string]*mock.Process
		EventServiceMap map[string]*mock.EventService
	}{
		ProcessMap:      processMap,
		EventServiceMap: eventServiceMap,
	}
}
