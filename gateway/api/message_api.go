package api

import (
	"github.com/it-chain/it-chain-Engine/gateway"
)

type MessageApi struct {
	grpcService gateway.GrpcService
}

func NewMessageApi(grpcService gateway.GrpcService) *ConnectionApi {
	return &ConnectionApi{
		grpcService: grpcService,
	}
}

//todo
func (c MessageApi) DeliverMessage(command gateway.MessageDeliverCommand) {

	//validation rule add
	c.grpcService.SendMessages(command.Body, command.Protocol, command.Recipients...)
}
