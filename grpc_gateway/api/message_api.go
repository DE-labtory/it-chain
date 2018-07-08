package api

import (
	"github.com/it-chain/it-chain-Engine/grpc_gateway"
)

type MessageApi struct {
	grpcService grpc_gateway.GrpcService
}

func NewMessageApi(grpcService grpc_gateway.GrpcService) *MessageApi {
	return &MessageApi{
		grpcService: grpcService,
	}
}

//todo validation rule added example( check length of recipent)
func (c MessageApi) DeliverMessage(command grpc_gateway.GrpcDeliverCommand) {

	//validation rule add
	c.grpcService.SendMessages(command.Body, command.Protocol, command.Recipients...)
}
