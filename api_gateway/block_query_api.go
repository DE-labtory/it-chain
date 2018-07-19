package api_gateway

import (
	"github.com/it-chain/engine/blockchain"
	"github.com/it-chain/yggdrasill"

	"sync"
)

type BlockQueryApi struct {
	CommitedBlockRepository
}

type CommitedBlockRepository interface {
	Save(block blockchain.Block) error
	GetLastBlock() (blockchain.Block, error)
	GetBlockByHeight(height blockchain.BlockHeight) (blockchain.Block, error)
}

type CommitedBlockRepositoryImpl struct {
	mux *sync.RWMutex
	yggdrasill.BlockStorageManager
}
