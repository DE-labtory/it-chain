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
	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/event"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/p2p"
	"github.com/it-chain/engine/p2p/api"
	"github.com/it-chain/engine/p2p/infra/adapter"
	"github.com/it-chain/engine/p2p/infra/mem"
	"github.com/it-chain/midgard"
)

type MockP2PProcess struct {
	id                 string
	eventHandler       adapter.EventHandler
	communicationApi   api.CommunicationApi
	leaderApi          api.LeaderApi
	grpcCommandHandler adapter.GrpcCommandHandler
	commandReceiver    chan interface{}
	eventReceiver      chan interface{}
	mockProcessTable   map[string]MockP2PProcess // map of mock process
}

func NewMockP2PProcess(mockProcessTable map[string]MockP2PProcess) MockP2PProcess {

	return MockP2PProcess{
		mockProcessTable: mockProcessTable,
	}
}

func (mpp *MockP2PProcess) Init(id string, ipAddress string) {

	peerRepository := mem.NewPeerReopository()

	//todo replace with component repo
	queryPeerRepository := api_gateway.PeerRepository{}

	peerQueryService := api_gateway.NewPeerQueryApi(&queryPeerRepository)

	communicationService := p2p.NewCommunicationService(mpp.Publish)

	communicationApi := api.NewCommunicationApi(&peerQueryService, communicationService)

	//todo should be replaced with real struct
	peerApi := &MockPeerApi{}

	//todo should be replaced with real struct
	eventService := common.EventServiceImpl{}

	mpp.eventHandler = adapter.NewEventHandler(&communicationApi, peerApi)

	leaderApi := api.NewLeaderApi(&peerRepository, &eventService)

	election := p2p.NewElection(0, "candidate", 0)

	electionRepository := p2p.NewElectionRepository(election)

	//todo should be replaced with mock client
	client := rpc.NewClient("")

	electionService := p2p.NewElectionService(electionRepository, &peerQueryService, client)

	pLTableService := p2p.PLTableServiceImpl{}

	mpp.grpcCommandHandler = adapter.NewGrpcCommandHandler(&leaderApi, electionService, &communicationApi, &pLTableService)

	mpp.communicationApi = api.NewCommunicationApi(&peerQueryService, communicationService)
}

func (mpp *MockP2PProcess) GetId() string {

	return mpp.id
}

//publish event => consumed by commandHandler and eventHandler through channel
func (mpp *MockP2PProcess) Publish(topic string, data interface{}) error {
	mpp.eventReceiver <- data

	return nil
}

func (mpp *MockP2PProcess) HandleEvent() {
	var done = true
	for done {

		e := <-mpp.eventReceiver

		switch reflect.TypeOf(e) {

		case reflect.TypeOf(event.ConnectionCreated{}):

			mpp.eventHandler.HandleConnCreatedEvent(e.(event.ConnectionCreated))

		case reflect.TypeOf(event.ConnectionClosed{}):

			mpp.eventHandler.HandleConnDisconnectedEvent(e.(event.ConnectionClosed))

		default:
			done = false

		}

	}
}

//send data to other mock process
func (mpp *MockP2PProcess) TriggerOutboundCommand(data interface{}, process MockP2PProcess) {
	process.commandReceiver <- data
}

//send data to other mock process
func (mpp *MockP2PProcess) TriggerOutboundEvent(data interface{}, process MockP2PProcess) {
	process.eventReceiver <- data
}

//running in go routine ex) go HandleCommand()
//deal with command
//
//if command is
func (mpp *MockP2PProcess) RpcListen(mpt map[string]MockP2PProcess) {
	var done = true
	for done {

		c := <-mpp.commandReceiver

		switch reflect.TypeOf(c) {

		case reflect.TypeOf(command.ReceiveGrpc{}):

			mpp.grpcCommandHandler.HandleMessageReceive(c.(command.ReceiveGrpc))

		case reflect.TypeOf(event.ConnectionCreated{}):
			connectionCreatedEvent := event.ConnectionCreated{
				ConnectionID: c.(event.ConnectionCreated).ConnectionID,
			}
			mpp.TriggerOutboundEvent(connectionCreatedEvent, mpt[connectionCreatedEvent.ConnectionID])

		case reflect.TypeOf(command.DeliverGrpc{}):
			switch c.(command.DeliverGrpc).Protocol {
			case "PLTableDeliverProtocol":

				// trigger grpc receive protocol to target peer

				grpcReceiveCommand := command.ReceiveGrpc{
					CommandModel: midgard.CommandModel{
						ID: c.(command.DeliverGrpc).ID,
					},
				}
				mpp.TriggerOutboundCommand(grpcReceiveCommand, mpt[grpcReceiveCommand.ID])

			}

		default:
			done = false

		}

	}
}

func (mpp *MockP2PProcess) RpcCall(data interface{}, process MockProcess) {

	mpp.commandReceiver <- data
}

func (mpp *MockP2PProcess) FindMockP2PProcess(id string) MockP2PProcess {

	return mpp.mockProcessTable[id]
}
