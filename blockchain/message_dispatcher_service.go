package blockchain

import "github.com/it-chain/yggdrasill/impl"

type MessageDispatcher interface {
	AddedBlock(block impl.DefaultBlock)
}
