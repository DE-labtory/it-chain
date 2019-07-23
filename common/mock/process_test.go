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
	"testing"

	"github.com/DE-labtory/it-chain/common/command"
	"github.com/DE-labtory/it-chain/common/mock"
	"github.com/magiconair/properties/assert"
)

func TestNewProcess(t *testing.T) {
	//when
	process := mock.NewProcess("1")

	//	then
	assert.Equal(t, process.Id, "1")

	//when
	process2 := mock.NewProcess("2")

	//	then
	assert.Equal(t, process2.Id, "2")

}

func TestProcess_Register(t *testing.T) {

	process := mock.NewProcess("1")

	test := &struct {
		Id string
	}{Id: "1"}

	process.Register(test)

	assert.Equal(t, process.Services, map[string]interface{}{"": test})
}

func TestProcess_RegisterHandler(t *testing.T) {

	process := mock.NewProcess("1")

	c := command.ReceiveGrpc{
		MessageId: "1",
	}

	process.RegisterHandler(func(command command.ReceiveGrpc) error {
		return nil
	})

	assert.Equal(t, process.GrpcCommandHandlers[0](c), nil)
}
