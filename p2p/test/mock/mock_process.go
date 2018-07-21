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
package mock

import (
	"reflect"

	"github.com/it-chain/engine/api_gateway"
	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/api"
	"github.com/it-chain/engine/p2p/infra/adapter"
	"github.com/it-chain/midgard"
)

type MockProcess interface {
	Init(id string, ipAddress string)

	GetId() string

	Publish(exchange string, topic string, data interface{}) error

	HandleEvent()

	HandleCommand()

	TriggerCommand(data interface{}, process MockProcess)

	TriggerEvent(data interface{}, process MockProcess)
}

type MockP2PProcess struct {
	id                 string
	eventHandler       adapter.EventHandler
	communicationApi   api.CommunicationApi
	leaderApi          api.LeaderApi
	grpcCommandHandler adapter.GrpcCommandHandler
	commandReceiver    <-chan interface{}
	eventReceiver      <-chan interface{}
	mockProcessTable   map[string]MockP2PProcess // map of mock process
}

func NewMockP2PProcess(mockProcessTable map[string]MockP2PProcess) MockP2PProcess {

	return MockP2PProcess{
		mockProcessTable: mockProcessTable,
	}
}

func (mpp *MockP2PProcess) Init(id string, ipAddress string) {

	peerRepository := api_gateway.NewPeerRepository(p2p.PLTable{
		Leader: p2p.Leader{
			LeaderId: p2p.LeaderId{
				Id: id,
			},
		},
		PeerTable: map[string]p2p.Peer{
			id: {
				PeerId: p2p.PeerId{
					Id: id,
				},
				IpAddress: ipAddress,
			},
		},
	})

	peerQueryService := api_gateway.NewPeerQueryApi(peerRepository)

	communicationService := adapter.NewCommunicationService(mpp.Publish)

	communicationApi := api.NewCommunicationApi(peerQueryService, communicationService)

	mpp.eventHandler = adapter.NewEventHandler(&communicationApi)

	leaderService := p2p.NewLeaderService()

	leaderApi := api.NewLeaderApi(leaderService, peerQueryService)

	election := p2p.NewElection(0, "candidate", 0)

	electionRepository := p2p.NewElectionRepository(election)

	electionService := p2p.NewElectionService(electionRepository, peerQueryService, mpp.Publish)

	pLTableService := p2p.NewPLTableService()

	mpp.grpcCommandHandler = adapter.NewGrpcCommandHandler(&leaderApi, electionService, &communicationApi, &pLTableService)

	mpp.communicationApi = api.NewCommunicationApi(peerQueryService, communicationService)
}

func (mpp *MockP2PProcess) GetId() string {
	return mpp.id
}

//publish event or command => consumed by commandHandler and eventHandler through channel
func (mpp *MockP2PProcess) Publish(exchange string, topic string, data interface{}) error {
	switch exchange {
	case "Command":
		mpp.commandReceiver <- data

	case "Event":
		mpp.eventReceiver <- data

	}

	return nil
}

func (mpp *MockP2PProcess) HandleEvent() {
	var done = true
	for done {

		event := <-mpp.eventReceiver

		switch reflect.TypeOf(event) {

		case reflect.TypeOf(p2p.ConnectionCreatedEvent{}):

			mpp.eventHandler.HandleConnCreatedEvent(event.(p2p.ConnectionCreatedEvent))

		case reflect.TypeOf(p2p.ConnectionDisconnectedEvent{}):

			mpp.eventHandler.HandleConnDisconnectedEvent(event.(p2p.ConnectionDisconnectedEvent))

		default:
			done = false

		}

	}
}

//send data to other mock process
func (mpp *MockP2PProcess) TriggerCommand(data interface{}, process MockP2PProcess) {
	process.commandReceiver <- data
}

//send data to other mock process
func (mpp *MockP2PProcess) TriggerEvent(data interface{}, process MockP2PProcess) {
	process.eventReceiver <- data
}

//running in go routine ex) go HandleCommand()
//deal with command
//
//if command is
func (mpp *MockP2PProcess) HandleCommand(mpt map[string]MockP2PProcess) {
	var done = true
	for done {

		command := <-mpp.commandReceiver

		switch reflect.TypeOf(command) {

		case reflect.TypeOf(p2p.GrpcReceiveCommand{}):

			mpp.grpcCommandHandler.HandleMessageReceive(command.(p2p.GrpcReceiveCommand))

		case reflect.TypeOf(p2p.ConnectionCreateCommand{}):
			connectionCreatedEvent := p2p.ConnectionCreatedEvent{
				EventModel: midgard.EventModel{
					ID: command.(p2p.ConnectionCreateCommand).ID,
				},
			}
			mpp.TriggerEvent(connectionCreatedEvent, mpt[connectionCreatedEvent.ID])

		case reflect.TypeOf(p2p.GrpcDeliverCommand{}):
			switch command.(p2p.GrpcDeliverCommand).Protocol {
			case "PLTableDeliverProtocol":

				// trigger grpc receive protocol to target peer

				grpcReceiveCommand := p2p.GrpcReceiveCommand{
					CommandModel: midgard.CommandModel{
						ID: command.(p2p.GrpcDeliverCommand).ID,
					},
				}
				mpp.TriggerCommand(grpcReceiveCommand, mpt[grpcReceiveCommand.ID])

			}

		default:
			done = false

		}

	}
}

func (mpp *MockP2PProcess) AddMockP2PProcess(mockProcess MockP2PProcess) {

	mpp.mockProcessTable[mockProcess.id] = mockProcess
}

func (mpp *MockP2PProcess) FindMockP2PProcess(id string) MockP2PProcess {

	return mpp.mockProcessTable[id]
}
