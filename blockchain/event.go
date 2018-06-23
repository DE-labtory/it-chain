package blockchain

import "github.com/it-chain/midgard"

type BlockValidatedEvent struct {
	midgard.EventModel
	Block Block
}
