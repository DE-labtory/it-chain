package comm

import (
	"google.golang.org/grpc"
	"it-chain/network/protos"
)

type handler func(message *message.Envelope)

type Connection struct{
	conn         *grpc.ClientConn
	cl           message.GossipClient
}