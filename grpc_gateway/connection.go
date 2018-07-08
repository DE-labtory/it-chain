package grpc_gateway

import "github.com/it-chain/midgard"

type Connection struct {
	midgard.AggregateModel
	Address string
}
