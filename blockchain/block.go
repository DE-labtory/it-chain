package blockchain

import (
	"github.com/it-chain/yggdrasill"
	"github.com/it-chain/yggdrasill/impl"
)

type Block = impl.DefaultBlock

type Repository interface {
	yggdrasill.BlockStorageManager
	NewEmptyBlock() (Block, error)
	GetBlockCreator() string
}
