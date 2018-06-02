package gateway

import (
	"github.com/it-chain/midgard"
)

//Connection 생성 command
type ConnectionCreateCommand struct {
	midgard.CommandModel
	Address string
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
}
