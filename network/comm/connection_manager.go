package comm

import (
	"it-chain/network/protos"
	"it-chain/service/domain"
)

type onError func(error)

//comm은 peer 들간의 connection을 유지하고있다.
//comm을 통해 peer들과 통신한다.
type ConnectionManager interface{

	SendStream(envelope message.Envelope, errorCallBack onError, connectionID string)

	Stop()

	Close(peerInfo domain.PeerInfo)

	CreateStreamConn(connectionID string, ip string, handle ReceiveMessageHandle) error

	Size() int
}