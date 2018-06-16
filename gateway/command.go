package gateway

import (
	"github.com/it-chain/it-chain-Engine/p2p"
	"github.com/it-chain/midgard"
)

//Connection 생성 command
type ConnectionCreateCommand struct {
	midgard.CommandModel
	Address string
}

//Connection close command
type ConnectionCloseCommand struct {
	midgard.CommandModel
}

//다른 Peer에게 Message전송 command
type MessageDeliverCommand struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}

//다른 Peer에게 Message수신 command
type MessageReceiveCommand struct {
	midgard.CommandModel
	Data         []byte
	ConnectionID string
	Protocol     string
	FromNode     p2p.Node // 메세지를 받는 경우 connectionID 만 있으면 수신자가 불편함이 있을 것 같아 추가했습니다.
	// 현재 node repo 에서 connection id 를 저장하지 않기 때문에 수신자가 nodeid 를 알 수 없습니다.
	// -frontalnh
}
