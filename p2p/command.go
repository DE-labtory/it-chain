package p2p

import "github.com/it-chain/midgard"

type GrpcRequestCommand struct {
	midgard.CommandModel
	Data         []byte
	ConnectionID string
	Protocol     string
}

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
