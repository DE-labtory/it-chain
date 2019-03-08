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

package mock

import (
	"reflect"
	"sync"

	"github.com/DE-labtory/engine/common/command"
)

type Process struct {
	mutex sync.Mutex
	Id    string

	GrpcCommandHandlers []func(command command.ReceiveGrpc) error
	GrpcCommandReceiver chan command.ReceiveGrpc //should be register to network's channel map
	Services            map[string]interface{}   // register service or api for testing which has injected mock client
}

func NewProcess(processId string) *Process {
	return &Process{
		Id:                  processId,
		Services:            make(map[string]interface{}),
		GrpcCommandReceiver: make(chan command.ReceiveGrpc),
		GrpcCommandHandlers: make([]func(command command.ReceiveGrpc) error, 0),
	}
}

// 테스트 과정에서 사용하기 위한 다양한 서비스 혹은 api 들을 등록한다.
// 이를 통해 특정 서비스 혹은 api를 사용함에 있어 어느 프로세스에 종속된 것인지를 직관적으로 알 수 있도록 한다.
// register all kinds of services or apis to use in test process
// so we can intuitively know processes we are using are belongs to the process
func (p *Process) Register(service interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.Services[reflect.ValueOf(service).Elem().Type().Name()] = service
}

// register handlers to be started when network manager start to work
func (p *Process) RegisterHandler(handler func(command command.ReceiveGrpc) error) error {

	p.GrpcCommandHandlers = append(p.GrpcCommandHandlers, handler)

	return nil
}
