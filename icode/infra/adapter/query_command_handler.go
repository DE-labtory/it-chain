package adapter

import (
	"github.com/it-chain/engine/common/amqp/rpc"
	"github.com/it-chain/engine/icode/api"
)

type QueryCommandHandler struct {
	server rpc.RpcServer
}

func NewQueryCommandHandler(server rpc.RpcServer) *QueryCommandHandler {
	return &QueryCommandHandler{
		server: server,
	}
}

func (q *QueryCommandHandler) SetHandleQueryCommand(api api.ICodeApi) error {
	return q.server.Register(api)
}
