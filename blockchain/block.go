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
	Dequeue() Block
}

type DefaultBlockPool struct {
	pool []Block
}

func (p *DefaultBlockPool) Enqueue(block Block) {
	p.pool = append(p.pool, block)
}

func (p *DefaultBlockPool) Dequeue() Block {
	if len(p.pool) == 0 {
		return nil
	}
	return p.pool[0]
}