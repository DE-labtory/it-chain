package api_gateway

import "github.com/it-chain/it-chain-Engine/blockchain"

type BlockQueryApi struct {
	BlockPoolRepository
}

type BlockPoolRepository interface {
	AddCreatedBlock(block blockchain.Block)
	GetStagedBlockByHeight(height blockchain.BlockHeight) (blockchain.Block, error)
	GetStagedBlockById(id string) (blockchain.Block, error)
}

type BlockPoolRepositoryImpl struct {
	Blocks []blockchain.Block
}
