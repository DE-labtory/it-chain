package p2p

import "github.com/it-chain/midgard"

type GrpcRequestCommand struct {
	midgard.CommandModel
	Data         []byte
	ConnectionID string
}
