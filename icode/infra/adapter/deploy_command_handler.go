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
	"os"

	"fmt"

	"github.com/it-chain/engine/common/command"
	"github.com/it-chain/engine/common/logger"
	"github.com/it-chain/engine/common/rabbitmq/rpc"
	"github.com/it-chain/engine/icode"
	"github.com/it-chain/engine/icode/api"
)

type DeployCommandHandler struct {
	icodeApi api.ICodeApi
}

func NewDeployCommandHandler(icodeApi api.ICodeApi) *DeployCommandHandler {
	return &DeployCommandHandler{
		icodeApi: icodeApi,
	}
}

func (d *DeployCommandHandler) HandleDeployCommand(deployCommand command.Deploy) (icode.Meta, rpc.Error) {

	savePath := os.Getenv("GOPATH") + "/src/github.com/it-chain/engine/.tmp/"
	logger.Info(nil, fmt.Sprintf("[Icode] icode deploying, url %s", deployCommand.Url))
	logger.Info(nil, fmt.Sprintf("[Icode] icode deploying, id %s", deployCommand.ID))
	logger.Info(nil, fmt.Sprintf("[Icode] icode saving path: %s", savePath))

	meta, err := d.icodeApi.Deploy(deployCommand.GetID(), savePath, deployCommand.Url, deployCommand.SshPath)

	if err != nil {
		logger.Error(nil, fmt.Sprintf("[Icode] fail to deploy icode, url %s", deployCommand.Url))
		return icode.Meta{}, rpc.Error{Message: err.Error()}
	}

	return meta, rpc.Error{}
}
