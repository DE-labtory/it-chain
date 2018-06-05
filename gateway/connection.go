package gateway

import "github.com/it-chain/midgard"

type Connection struct {
	midgard.AggregateModel
	Address string
}
