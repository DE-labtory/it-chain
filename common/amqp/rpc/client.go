package infra

import (
	"net/rpc"

	"github.com/streadway/amqp"
	"github.com/vibhavp/amqp-rpc"
)

type RpcClient struct {
	conn        amqp.Connection
	clientCodec *rpc.ClientCodec
	client      *rpc.Client
}

func NewRpcClient(conn amqp.Connection, routingKey string) (*RpcClient, error) {
	codec, err := amqprpc.NewClientCodec(&conn, routingKey, amqprpc.JSONCodec{})
	if err != nil {
		return nil, err
	}

	return &RpcClient{
		conn:        conn,
		clientCodec: &codec,
		client:      rpc.NewClientWithCodec(codec),
	}, nil
}

func (r *RpcClient) Call(serviceMethod string, args interface{}, reply interface{}) error {
	return r.client.Call(serviceMethod, args, reply)
}
