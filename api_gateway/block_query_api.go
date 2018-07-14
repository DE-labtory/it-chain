package api_gateway

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/blockchain"
)

var ErrNoStagedBlock = errors.New("Error can not find staged block")

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

type BlockPoolRepositoryImpl struct {
	Blocks []blockchain.Block
}

func NewBlockpoolRepositoryImpl() *BlockPoolRepositoryImpl {
	return &BlockPoolRepositoryImpl{
		Blocks: make([]blockchain.Block, 0),
	}
}

func (bpr *BlockPoolRepositoryImpl) GetStagedBlockByHeight(blockHeight uint64) (blockchain.Block, error) {
	for _, block := range bpr.Blocks {
		if block.GetHeight() == blockHeight {
			return block, nil
		}
	}
	return nil, ErrNoStagedBlock
}
func (bpr *BlockPoolRepositoryImpl) GetStagedBlockById(blockId string) (blockchain.Block, error) {

}
