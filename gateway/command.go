package gateway

import (
	"github.com/it-chain/midgard"
)

type ConnectionCreateCommand struct {
	midgard.CommandModel
	Address string
}

type MessageDeliverCommand struct {
	midgard.CommandModel
	Recipients []string
	Body       []byte
	Protocol   string
}

type GrpcRequestCommand struct {
	midgard.CommandModel
	Data         []byte
	ConnectionID string
}
