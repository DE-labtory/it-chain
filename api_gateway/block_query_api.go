package api_gateway

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
)

type BlockQueryApi struct {
	blockpoolRepository     BlockPoolRepository
	commitedBlockRepository CommitedBlockRepository
}

func (b BlockQueryApi) GetStagedBlockByHeight(height uint64) (blockchain.Block, error) {
	return b.blockpoolRepository.GetStagedBlockByHeight(height)
}

func (b BlockQueryApi) GetStagedBlockById(blockId string) (blockchain.Block, error) {
	return b.blockpoolRepository.GetStagedBlockById(blockId)
}

func (b BlockQueryApi) GetLastCommitedBlock() (blockchain.Block, error) {
	return b.commitedBlockRepository.GetLastBlock()
}

func (b BlockQueryApi) GetCommitedBlockByHeight(height uint64) (blockchain.Block, error) {
	return b.commitedBlockRepository.GetBlockByHeight(height)
}

type BlockPoolRepository interface {
	GetStagedBlockByHeight(blockHeight uint64) (blockchain.Block, error)
	GetStagedBlockById(blockId string) (blockchain.Block, error)
}

type CommitedBlockRepository interface {
	GetLastBlock() (blockchain.Block, error)
	GetBlockByHeight(height uint64) (blockchain.Block, error)
}
