package blockchain

import "github.com/it-chain/midgard"

type BlockCreatedEvent struct {
	midgard.EventModel
	Block Block
}
