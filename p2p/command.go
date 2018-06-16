package p2p

import "github.com/it-chain/midgard"

type GrpcRequestCommand struct {
	midgard.CommandModel
	Data         []byte
	ConnectionID string
	Protocol     string
	FromNode     Node
	ToNode       Node
}

type MessageDeliverCommand struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}

//다른 Peer에게 Message수신 command =>
// gate way의 MessageReceiverCommand와 정확히 동일 한 값을 가져야 하므로 사용하지 않았습니다.
// 훗날 gateway의 MessageReceiverCommand가 변경되는 경우 동일한 로직으로 처리되어야 합니다.
//type MessageReceiveCommand struct {
//	midgard.CommandModel
//	Data         []byte
//	ConnectionID string
//	Protocol     string
//}
