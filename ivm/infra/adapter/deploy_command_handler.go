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
	"fmt"
	"os"

	"github.com/it-chain/engine/common"
	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/ivm"
	"github.com/it-chain/engine/ivm/api"
)

type DeployCommandHandler struct {
	icodeApi api.ICodeApi
}

func NewDeployCommandHandler(icodeApi api.ICodeApi) *DeployCommandHandler {
	return &DeployCommandHandler{
		icodeApi: icodeApi,
	}
}

func (d *DeployCommandHandler) HandleDeployCommand(deployCommand command.Deploy) (ivm.ICode, rpc.Error) {

	savePath := os.Getenv("GOPATH") + "/src/github.com/it-chain/engine/.tmp/"

	absolutePath, err := common.RelativeToAbsolutePath(deployCommand.SshPath)
	if err != nil {
		logger.Error(nil, fmt.Sprintf("[Icode] fail to get ssh path : %s", absolutePath))
		return ivm.ICode{}, rpc.Error{Message: err.Error()}
	}
	icode, err := d.icodeApi.Deploy(deployCommand.ICodeId, savePath, deployCommand.Url, absolutePath, "")

	if err != nil {
		logger.Error(nil, fmt.Sprintf("[Icode] fail to deploy ivm, url %s", deployCommand.Url))
		return ivm.ICode{}, rpc.Error{Message: err.Error()}
	}

	return icode, rpc.Error{}
}
