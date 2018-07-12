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

func (c MessageApi) DeliverMessage(body []byte, protocol string, ids ...string) {

	//validation rule add
	c.grpcService.SendMessages(body, protocol, ids...)
}
