package blockchain

import (
	"github.com/it-chain/midgard"
	"github.com/it-chain/yggdrasill/impl"
)

type BlockAddCommand struct {
	midgard.CommandModel
	impl.DefaultBlock
}
