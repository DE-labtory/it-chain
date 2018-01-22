package comm

import (
	"it-chain/network/protos"
	"it-chain/service/domain"
)

type onError func(error)

//comm은 peer 들간의 connection을 유지하고있다.
//comm을 통해 peer들과 통신한다.
type Comm interface{

	SendStream(envelope message.Envelope, errorCallBack onError, connectionID string)

	Stop()

	Close(peerInfo domain.PeerInfo)

	//connection이 유지
	CreateStreamConn(connectionID string, ip string) error

	//connection이 유지되지 않는다.
	Send(ip string) error

	Size() int
}