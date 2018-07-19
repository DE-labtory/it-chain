package adapter

import (
	"fmt"
	"log"

	"github.com/it-chain/engine/icode"
	"github.com/it-chain/engine/icode/api"
)

type UnDeployCommandHandler struct {
	icodeApi api.ICodeApi
}

func NewUnDeployCommandHandler(icodeApi api.ICodeApi) *UnDeployCommandHandler {
	return &UnDeployCommandHandler{
		icodeApi: icodeApi,
	}
}

func (u *UnDeployCommandHandler) HandleUnDeployCommand(command icode.UnDeployCommand) {
	err := u.icodeApi.UnDeploy(command.ID)
	if err != nil {
		log.Println(fmt.Sprintf("error in handle undeploy command %s", err.Error()))
	}
}
