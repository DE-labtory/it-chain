package adapter

import "github.com/it-chain/engine/icode/api"

type CommandHandler struct {
	iCodeApi api.ICodeApi
}

func NewCommandHandler(icodeApi api.ICodeApi) *CommandHandler {

	return &CommandHandler{
		iCodeApi: icodeApi,
	}
}
