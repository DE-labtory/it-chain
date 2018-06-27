package p2p

import "github.com/it-chain/midgard"

type GrpcRequestCommand struct {
	midgard.CommandModel
	Data         []byte
	ConnectionID string
	Protocol     string
	FromPeer     Peer
	ToPeer       Peer
}

type MessageDeliverCommand struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}