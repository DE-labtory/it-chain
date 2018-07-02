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
	Delete(height BlockHeight)
}

type BlockPoolModel struct {
	pool map[BlockHeight] Block
}

func (p *BlockPoolModel) Add(block Block) {
	height := block.GetHeight()
	p.pool[height] = block
}

func (p *BlockPoolModel) Get(height BlockHeight) Block {
	return p.pool[height]
}

func (p *BlockPoolModel) Delete(height BlockHeight) {
	delete(p.pool, height)
}