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

package adapter

import (
	"fmt"
	"os"

	"github.com/DE-labtory/it-chain/common/command"
	"github.com/DE-labtory/it-chain/common/rabbitmq/rpc"
	"github.com/DE-labtory/it-chain/ivm"
	"github.com/DE-labtory/it-chain/ivm/api"
	"github.com/DE-labtory/iLogger"
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
	savePath := os.Getenv("GOPATH") + "/src/github.com/DE-labtory/it-chain/.tmp/"

	if deployCommand.SshPath != "" {
		icode, err := d.icodeApi.Deploy(savePath, deployCommand.Url, deployCommand.SshPath, deployCommand.Password)
		if err != nil {
			iLogger.Error(nil, fmt.Sprintf("[Icode] fail to deploy ivm, url %s", deployCommand.Url))
			return ivm.ICode{}, rpc.Error{Message: err.Error()}
		}

		return icode, rpc.Error{}
	} else if len(deployCommand.SshRaw) != 0 {
		icode, err := d.icodeApi.DeployFromRawSsh(savePath, deployCommand.Url, deployCommand.SshRaw, deployCommand.Password)
		if err != nil {
			iLogger.Error(nil, fmt.Sprintf("[Icode] fail to deploy ivm, url %s", deployCommand.Url))
			return ivm.ICode{}, rpc.Error{Message: err.Error()}
		}
		return icode, rpc.Error{}
	} else {
		icode, err := d.icodeApi.Deploy(savePath, deployCommand.Url, "", deployCommand.Password)
		if err != nil {
			iLogger.Error(nil, fmt.Sprintf("[Icode] fail to deploy ivm, url %s", deployCommand.Url))
			return ivm.ICode{}, rpc.Error{Message: err.Error()}
		}
		return icode, rpc.Error{}
	}

}
