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

package adapter

import (
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/icode"
	"github.com/it-chain/engine/icode/api"
)

type IcodeExecuteCommandHandler struct {
	iCodeApi api.ICodeApi
}

func NewIcodeExecuteCommandHandler(icodeApi api.ICodeApi) IcodeExecuteCommandHandler {
	return IcodeExecuteCommandHandler{
		iCodeApi: icodeApi,
	}
}

func (i *IcodeExecuteCommandHandler) HandleTransactionExecuteCommandHandler(command command.ExecuteICode) (icode.Result, rpc.Error) {

	request := icode.Request{
		Args:     command.Args,
		Function: command.Function,
		ICodeID:  command.ID,
		Type:     command.Method,
	}

	result, err := i.iCodeApi.ExecuteRequest(request)

	if err != nil {
		return icode.Result{}, rpc.Error{Message: err.Error()}
	}

	return result, rpc.Error{}
}
