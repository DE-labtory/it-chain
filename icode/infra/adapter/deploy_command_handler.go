package adapter

import (
	"fmt"
	"log"

	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/icode"
	"github.com/it-chain/it-chain-Engine/icode/api"
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

	_, err := d.icodeApi.Deploy(conf.GetConfiguration().Icode.ICodeSavePath, command.Url, command.SshPath)
	if err != nil {
		log.Println(fmt.Sprintf("error in handle deploy command %s", err.Error()))
	}
}
