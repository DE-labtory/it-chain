package api_gateway

import "github.com/it-chain/it-chain-Engine/blockchain"

type BlockQueryApi struct {
	createdBlockRepository  CreatedBlockRepository
	stagedBlockRepository   StagedBlockRepository
	commitedBlockRepository CommitedBlockRepository
}

func (b BlockQueryApi) GetStagedBlockByHeight(height uint64) (blockchain.Block, error) {
	return b.stagedBlockRepository.GetBlockByHeight(height)
}

func (b BlockQueryApi) GetStagedBlockById(blockId string) (blockchain.Block, error) {
	return b.stagedBlockRepository.GetBlockById(blockId)
}

func (b BlockQueryApi) GetLastCommitedBlock() (blockchain.Block, error) {
	return b.commitedBlockRepository.GetLastBlock()
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
