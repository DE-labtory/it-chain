package gateway

import (
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
type GrpcDeliverCommand struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}

//다른 Peer에게 Message수신 command
type GrpcReceiveCommand struct {
	midgard.CommandModel
	Data         []byte
	ConnectionID string
	Protocol     string
}
