package adapter

import (
	"fmt"
	"log"

	"github.com/it-chain/engine/conf"
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

func (d *DeployCommandHandler) HandleDeployCommand(command icode.DeployCommand) {
	command.GetID()
	_, err := d.icodeApi.Deploy(command.GetID(), conf.GetConfiguration().Icode.ICodeSavePath, command.Url, command.SshPath)
	if err != nil {
		log.Println(fmt.Sprintf("error in handle deploy command %s", err.Error()))
	}
}
