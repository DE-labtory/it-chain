package blockchain

import (
	"github.com/it-chain/yggdrasill"
	"github.com/it-chain/yggdrasill/common"
	"github.com/it-chain/yggdrasill/impl"
)

type Block = common.Block

type DefaultBlock = impl.DefaultBlock

type Repository interface {
	yggdrasill.BlockStorageManager
	NewEmptyBlock() (Block, error)
	GetBlockCreator() string
}

type BlockPool interface {
	Enqueue(block Block)
	Dequeue() error
}