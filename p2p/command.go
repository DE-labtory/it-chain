package p2p

import "github.com/it-chain/midgard"

//다른 Peer에게 Message수신 command
type GrpcReceiveCommand struct {
	midgard.CommandModel
	Body         []byte
	ConnectionID string
	Protocol     string
}

//다른 Peer에게 Message전송 command
type GrpcDeliverCommand struct {
	midgard.CommandModel
	Recipients []string //connectionId
	Body       []byte
	Protocol   string
}