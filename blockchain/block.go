package blockchain

import (
	"github.com/it-chain/yggdrasill"
	"github.com/it-chain/yggdrasill/common"
	"github.com/it-chain/yggdrasill/impl"
)

type Block = common.Block

type Repository interface {
	yggdrasill.BlockStorageManager
	NewEmptyBlock() (*impl.DefaultBlock, error)
	GetBlockCreator() string
}
