package adapter

import (
	"github.com/it-chain/it-chain-Engine/gateway"
	"github.com/it-chain/it-chain-Engine/gateway/api"
)

type ConnectionCommandHandler struct {
	connectionApi api.ConnectionApi
}

type MessageCommandHandler struct {
	messageApi api.MessageApi
}

func (c ConnectionCommandHandler) HandleGrpcDeliverCommand(command gateway.GrpcDeliverCommand) {

}
