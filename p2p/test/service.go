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

type ProcessIdentity struct {
	Id        string
	IpAddress string
}

//
//func SetTestEnvironment(processList []ProcessIdentity) map[string]*mock.Process {
//	networkManager := mock.NewNetworkManager()
//
//	m := make(map[string]*mock.Process)
//	for _, entity := range processList {
//		// avengers - create process and network
//		process := mock.NewProcess()
//		networkManager.AddProcess(process)
//
//		election := p2p.NewElection(entity.Id, 30, p2p.Ticking, 0)
//		peerRepository := mem.NewPeerReopository()
//		savePeerList(peerRepository, processList)
//
//		peerQueryService := api_gateway.NewPeerQueryApi(&peerRepository)
//
//		// avengers - mock client, server
//		client := mock.NewClient(entity.Id, networkManager.GrpcCall)
//		server := mock.NewServer(entity.Id, networkManager.GrpcConsume)
//
//		eventService := mock2.MockEventService{}
//		eventService.PublishFunc = func(topic string, event interface{}) error {
//			return nil
//		}
//
//		// inject avengers client
//		electionService := p2p.NewElectionService(&election, &peerQueryService, &client)
//
//		pLTableService := p2p.PLTableService{}
//
//		// inject avengers client
//		communicationService := p2p.NewCommunicationService(&client)
//
//		communicationApi := api.NewCommunicationApi(&peerQueryService, communicationService)
//
//		leaderApi := api.NewLeaderApi(&peerRepository, &eventService)
//		grpcCommandHandler := adapter.NewGrpcCommandHandler(&leaderApi, &electionService, &communicationApi, pLTableService)
//
//		// avengers server register command handler
//		server.Register("message.receive", grpcCommandHandler.HandleMessageReceive)
//
//		process.Register(&electionService)
//		process.Register(&peerRepository)
//		m[process.Id] = &process
//	}
//
//	logger.Infof(nil, "created process: %v", m)
//	return m
//}
