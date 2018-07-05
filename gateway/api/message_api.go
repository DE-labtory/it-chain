package api

import (
	"github.com/it-chain/it-chain-Engine/gateway"
)

type MessageApi struct {
	grpcService gateway.GrpcService
}

func NewMessageApi(grpcService gateway.GrpcService) *MessageApi {
	return &MessageApi{
		grpcService: grpcService,
	}
}

//todo validation rule added example(check length of recipent)
func (c MessageApi) DeliverMessage(body []byte, protocol string, ids ...string) {

	//validation rule add
	c.grpcService.SendMessages(body, protocol, ids...)
}
