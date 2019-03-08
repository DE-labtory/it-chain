/*
 * Copyright 2018 DE-labtory
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

package mock_test

import (
	"sync"
	"testing"
	"time"

	"github.com/DE-labtory/engine/common/command"
	"github.com/DE-labtory/engine/common/mock"
	"github.com/magiconair/properties/assert"
)

func TestNewNetworkManager(t *testing.T) {
}

func TestNetworkManager_AddProcess(t *testing.T) {

}

func TestNetworkManager_GrpcCall(t *testing.T) {
	tests := map[string]struct {
		input struct {
			RecipientList []string
			Protocol      string
		}
	}{
		"success": {input: struct {
			RecipientList []string
			Protocol      string
		}{RecipientList: []string{"1", "2"}, Protocol: "test"}},
		"no receiver test": {input: struct {
			RecipientList []string
			Protocol      string
		}{RecipientList: []string{}, Protocol: "test"}},
	}

	for testName, test := range tests {
		networkManager := mock.NewNetworkManager()
		t.Logf("running test case %s", testName)

		deliverGrpc := &command.DeliverGrpc{
			RecipientList: test.input.RecipientList,
			Protocol:      test.input.Protocol,
		}
		networkManager.GrpcCall("1", "message.deliver", *deliverGrpc, func() {})

		t.Logf("end of test calling")

		for _, processId := range test.input.RecipientList {
			t.Logf("processId:%s is receiving", processId)

			go func(processId string) {
				a := <-networkManager.ChannelMap[processId]["message.receive"]
				assert.Equal(t, a.Protocol, test.input.Protocol)
			}(processId)
		}
	}
}

//todo error
//func TestNetworkManager_GrpcConsume(t *testing.T) {
//	callbackIndex := 1
//
//	tests := map[string]struct {
//		input struct {
//			RecipientList []string
//			ProcessId     string
//			handler       func(c command.ReceiveGrpc) error
//		}
//	}{
//		"success": {input: struct {
//			RecipientList []string
//			ProcessId     string
//			handler       func(c command.ReceiveGrpc) error
//		}{RecipientList: []string{"1", "2", "3"},
//			ProcessId: "1",
//			handler:   func(c command.ReceiveGrpc) error { callbackIndex = 2; t.Logf("handler!"); return nil }}},
//	}
//
//	for testName, test := range tests {
//		t.Logf("running test case %s", testName)
//
//		networkManager := mock.NewNetworkManager()
//
//		deliverGrpc := &command.DeliverGrpc{
//			RecipientList: test.input.RecipientList,
//		}
//		networkManager.ChannelMap[test.input.ProcessId] = make(map[string]chan command.ReceiveGrpc)
//		networkManager.ChannelMap[test.input.ProcessId]["message.receive"] = make(chan command.ReceiveGrpc)
//
//		networkManager.GrpcConsume(test.input.ProcessId, "message.receive", test.input.handler)
//
//		networkManager.GrpcCall("1", "message.deliver", *deliverGrpc, func() {
//			callbackIndex++
//		})
//
//		t.Logf("end of calling!")
//
//	}
//
//	time.Sleep(5 * time.Second)
//	assert.Equal(t, callbackIndex, 2)
//}

func TestNetworkManager_Publish(t *testing.T) {
	mem := ClosureMemory()
	wg := sync.WaitGroup{}
	wg.Add(2)
	networkManager := SetNetworkManager(mem)

	event := command.DeliverGrpc{
		MessageId:     "1",
		RecipientList: []string{"2", "3"},
	}

	go func() {
		m2 := <-networkManager.ProcessMap["2"].GrpcCommandReceiver
		assert.Equal(t, m2.MessageId, "1")
		wg.Done()
	}()

	go func() {
		m3 := <-networkManager.ProcessMap["3"].GrpcCommandReceiver
		assert.Equal(t, m3.MessageId, "1")
		wg.Done()
	}()

	networkManager.Publish("1", "deliver.message", event)

	wg.Wait()
}

func TestNetworkManager_Start(t *testing.T) {
	mem := ClosureMemory()
	net := SetNetworkManager(mem)

	net.Start()

	command := command.DeliverGrpc{
		MessageId:     "123",
		RecipientList: []string{"2", "3"},
	}

	net.Publish("1", "message.deliver", command)
	net.Publish("1", "message.deliver", command)

	time.Sleep(2 * time.Second)

	assert.Equal(t, mem(), 5)
}

func SetNetworkManager(closerMemory func() int) *mock.NetworkManager {
	networkManager := mock.NewNetworkManager()

	process1 := mock.NewProcess("1")
	handler1 := func(command command.ReceiveGrpc) error {
		return nil
	}
	process1.RegisterHandler(handler1)

	process2 := mock.NewProcess("2")
	handler2 := func(command command.ReceiveGrpc) error {
		closerMemory()
		return nil
	}
	process2.RegisterHandler(handler2)

	process3 := mock.NewProcess("3")
	handler3 := func(command command.ReceiveGrpc) error {
		closerMemory()
		return nil
	}
	process3.RegisterHandler(handler3)

	networkManager.AddProcess(process1)
	networkManager.AddProcess(process2)
	networkManager.AddProcess(process3)

	return networkManager
}

func ClosureMemory() func() int {
	i := 0

	return func() int {
		i++
		return i
	}
}
