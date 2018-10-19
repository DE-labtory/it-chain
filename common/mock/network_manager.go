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
	"sync"
	"time"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/iLogger"
)

//network manager builds environment for communication between multiple nodes in network

type NetworkManager struct {
	mutex      sync.Mutex
	ChannelMap map[string]map[string]chan command.ReceiveGrpc // channel for receive deliverGrpc command

	ProcessMap map[string]*Process
}

func NewNetworkManager() *NetworkManager {

	return &NetworkManager{
		ChannelMap: make(map[string]map[string]chan command.ReceiveGrpc),
		ProcessMap: make(map[string]*Process),
		mutex:      sync.Mutex{},
	}
}

//receiver => processId
//queue name => queue
func (n *NetworkManager) Push(processId string, queue string, c command.ReceiveGrpc) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	iLogger.Infof(nil, "push to channel: %s", processId)
	n.ChannelMap[processId][queue] <- c

	return nil
}

//Grpc call function would be injected to rpc client
//processId => sender
//params.RecipientList => receiver
func (n *NetworkManager) GrpcCall(processId string, queue string, params interface{}, callback interface{}) error {

	//find receiver process and deliver command through channel
	for _, v := range params.(command.DeliverGrpc).RecipientList {
		//convert grpc deliver message to grpc receive message
		extracted := command.ReceiveGrpc{
			Body:         params.(command.DeliverGrpc).Body,
			Protocol:     params.(command.DeliverGrpc).Protocol,
			ConnectionID: processId, // sender of this message
		}

		go func(v string, queue string) {
			queue = "message.receive"
			n.Push(v, queue, extracted)

		}(v, queue)
	}

	return nil
}

// Will Be DEPRECATED!
//GrpcConsume would be injected to rpc server
//processId => receiver
func (n *NetworkManager) GrpcConsume(processId string, queue string, handler func(command command.ReceiveGrpc) error) error {

	//start command distributer
	go func(processId string, queue string) {

		end := true

		for end {
			select {
			case message := <-n.ChannelMap[processId][queue]:
				iLogger.Infof(nil, "receive message from : %s message: %v", processId, message)
				handler(message)

			case <-time.After(4 * time.Second):
				iLogger.Info(nil, "failed to consume, timed out!")
				end = false
			}
		}
	}(processId, queue)

	return nil
}

// test success
func (n *NetworkManager) Publish(from string, topic string, event interface{}) error {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	if reflect.ValueOf(event).Type().Name() == "DeliverGrpc" {
		for _, id := range event.(command.DeliverGrpc).RecipientList {
			iLogger.Infof(nil, "publish to: %s", id)

			go func(id string) {
				n.ProcessMap[id].GrpcCommandReceiver <- command.ReceiveGrpc{
					MessageId:    event.(command.DeliverGrpc).MessageId,
					Body:         event.(command.DeliverGrpc).Body,
					ConnectionID: from,
					Protocol:     event.(command.DeliverGrpc).Protocol,
				}
			}(id)

		}
	}

	return nil
}

func (n *NetworkManager) Start() {
	for id, process := range n.ProcessMap {
		go func(id string, process *Process) {
			iLogger.Infof(nil, "process %s is running", process.Id)
			end := true
			for end {
				select {
				case message := <-process.GrpcCommandReceiver:
					iLogger.Infof(nil, "receive message from : %s message: %v", id, message)
					for _, handler := range process.GrpcCommandHandlers {
						handler(message)
					}

				case <-time.After(4 * time.Second):
					iLogger.Info(nil, "failed to consume, timed out!")
					end = false
				}
			}
		}(id, process)
	}
}

func (n *NetworkManager) AddProcess(process *Process) {
	n.mutex.Lock()
	defer n.mutex.Unlock()

	n.ProcessMap[process.Id] = process
}
