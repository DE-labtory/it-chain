package blockchain

import "github.com/it-chain/yggdrasill/impl"

type BlockRequestMessage struct {
	Height uint64
}

type BlockResponseMessage struct {
	Block *impl.DefaultBlock
}
