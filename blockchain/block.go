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

type BlockHeight = uint64

type BlockPool interface {
	Add(block Block)
	Get(height BlockHeight) Block
}

type DefaultBlockPool struct {
	pool map[BlockHeight] Block
}

func (p *DefaultBlockPool) Add(block Block) {
	height := block.GetHeight()
	p.pool[height] = block
}

func (p *DefaultBlockPool) Get(height BlockHeight) Block {
	return p.pool[height]
}