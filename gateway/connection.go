package gateway

import "github.com/it-chain/midgard"

type Connection struct {
	midgard.Aggregate
	ID      string
	Address string
}
