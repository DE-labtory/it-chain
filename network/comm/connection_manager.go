package comm

import (
	pb "it-chain/network/protos"
	"it-chain/network/comm/msg"
)

type OnError func(error)
type OnSuccess func(interface{})

//comm은 peer 들간의 connection을 유지하고있다.
//comm을 통해 peer들과 통신한다.
type ConnectionManager interface{

	SendStream(data *pb.StreamMessage, successCallBack OnSuccess, errorCallBack OnError, connectionID string)

	Stop()

	Close(connectionID string)

	CreateStreamClientConn(connectionID string, ip string) error

	Size() int

	//Server on function
	Stream(stream pb.StreamService_StreamServer) (error)

	SetOnConnectHandler(onConnectionHandler OnConnectionHandler)

	Subscribe(name string, subfunc func(message msg.OutterMessage))
}