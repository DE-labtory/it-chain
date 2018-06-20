package blockchain

import (
	"github.com/it-chain/yggdrasill/impl"
)

type BlockRequestMessage struct {
	TimeUnix int64
}

type BlockResponseMessage struct {
	BlockInfo impl.DefaultBlock
}
