package api_gateway

import (
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/yggdrasill"

	"sync"
)

type BlockQueryApi struct {
	BlockPoolRepository
	CommitedBlockRepository
}

type BlockPoolRepository interface {
	AddCreatedBlock(block blockchain.DefaultBlock)
	GetStagedBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
	GetStagedBlockById(id string) (blockchain.DefaultBlock, error)
	GetFirstStagedBlock() (blockchain.DefaultBlock, error)
}

type BlockPoolRepositoryImpl struct {
	Blocks []blockchain.Block
}

type CommitedBlockRepository interface {
	Save(block blockchain.DefaultBlock) error
	GetLastBlock() (blockchain.DefaultBlock, error)
	GetBlockByHeight(height blockchain.BlockHeight) (blockchain.DefaultBlock, error)
}

type CommitedBlockRepositoryImpl struct {
	mux *sync.RWMutex
	yggdrasill.BlockStorageManager
}
