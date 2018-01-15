package comm

import (
	"it-chain/network/protos"
	"it-chain/service/peer"
)

//comm은 peer 들간의 connection을 유지하고있다.
//comm을 통해 peer들과 통신한다.
type Comm interface{

	Send(envelop message.Envelope, peerInfos []peer.PeerInfo)

	Stop()

	Close(peerInfo peer.PeerInfo)

	CloseAll()

	CreateConn(peerInfo peer.PeerInfo) error
}