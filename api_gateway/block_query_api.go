package api_gateway

import "github.com/it-chain/it-chain-Engine/blockchain"

type BlockQueryApi struct {
	createdBlockRepository  CreatedBlockRepository
	stagedBlockRepository   StagedBlockRepository
	commitedBlockRepository CommitedBlockRepository
}

type CreatedBlockRepository interface {
}

type StagedBlockRepository interface {
	GetBlockByHeight(blockHeight uint64) (blockchain.Block, error)
	GetBlockById(blockId string) (blockchain.Block, error)
}

type CommitedBlockRepository interface {
	GetLastBlock() (blockchain.Block, error)
}
