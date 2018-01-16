package comm

import (
	"it-chain/network/protos"
	"it-chain/service/domain"
)

type onError func(error)

//comm은 peer 들간의 connection을 유지하고있다.
//comm을 통해 peer들과 통신한다.
type Comm interface{

	Send(envelop message.Envelope, errorCallBack onError, peerInfos ...domain.PeerInfo)

	Stop()

	Close(peerInfo domain.PeerInfo)

	CreateConn(peerInfo domain.PeerInfo) error

	Size() int
}