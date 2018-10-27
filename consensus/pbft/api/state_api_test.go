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

package api_test

import (
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/it-chain/engine/common/mock"
	"github.com/it-chain/engine/consensus/pbft"
	"github.com/it-chain/engine/consensus/pbft/api"
	"github.com/it-chain/engine/consensus/pbft/infra/adapter"
	"github.com/it-chain/engine/consensus/pbft/infra/mem"
	"github.com/it-chain/iLogger"
	"github.com/stretchr/testify/assert"
)

func TestStateApi_Integration_Test(t *testing.T) {

	tests := map[string]struct {
		input struct {
			processList   []string
			consensusTime int
		}
	}{
		"4 node test": {
			input: struct {
				processList   []string
				consensusTime int
			}{
				processList:   generateProcessList(4),
				consensusTime: 20,
			},
		},
		"8 node test": {
			input: struct {
				processList   []string
				consensusTime int
			}{
				processList:   generateProcessList(8),
				consensusTime: 20,
			},
		},
		"16 node test": {
			input: struct {
				processList   []string
				consensusTime int
			}{
				processList:   generateProcessList(16),
				consensusTime: 20,
			},
		},
	}

	for testName, test := range tests {
		t.Logf("running test case %s", testName)

		env := SetTestEnvironment(test.input.processList)

		for _, p := range env.ProcessMap {
			process := *p
			pRepo := process.Services["ParliamentRepository"].(*mem.ParliamentRepository)
			p := pRepo.Load()
			p.SetLeader("1")
			pRepo.Save(p)
		}

		stateApi1 := env.ProcessMap["1"].Services["StateApi"].(*api.StateApi)
		stateRepos := make([]*mem.StateRepository, 0)
		states := make([]pbft.State, 0)

		for _, p := range test.input.processList {
			stateRepo := env.ProcessMap[p].Services["StateRepository"].(*mem.StateRepository)
			tempState, _ := stateRepo.Load()
			states = append(states, tempState)
			stateRepos = append(stateRepos, stateRepo)
		}

		proposedBlock := pbft.ProposedBlock{Seal: []byte{'s', 'd', 'f'}, Body: []byte{'2', '3', '3'}}

		stateApi1.StartConsensus(proposedBlock)
		time.Sleep(time.Duration(test.input.consensusTime * int(time.Second)))

		for i := 0; i < len(stateRepos)-1; i++ {
			iLogger.Infof(nil, "SEAL: [%s]", states[i].Block.Seal)
			assert.Equal(t, states[i].Block.Seal, states[i+1].Block.Seal)
		}

	}
}

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

		randSecond := rand.Intn(10) * int(time.Second)
		eventService.SetDelayTime(time.Duration(randSecond))

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

	iLogger.Infof(nil, "[consensus] created process: %v", processMap)
	iLogger.Infof(nil, "[consensus] created event service: %v", eventServiceMap)

	networkManager.Start()
	return struct {
		ProcessMap      map[string]*mock.Process
		EventServiceMap map[string]*mock.EventService
	}{
		ProcessMap:      processMap,
		EventServiceMap: eventServiceMap,
	}
}

func generateProcessList(processNum int) []string {
	processList := make([]string, 0)

	for i := 0; i < processNum; i++ {
		processList = append(processList, strconv.Itoa(i+1))
	}

	return processList
}
