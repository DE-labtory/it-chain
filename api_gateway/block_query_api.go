package api_gateway

import (
	"errors"
	"sync"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/yggdrasill"
)

var ErrNoStagedBlock = errors.New("Error can not find staged block")
var ErrGetCommitedBlock = errors.New("Error in getting commited block")

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
	AddCreatedBlock(block blockchain.Block)
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

func (bpr *BlockPoolRepositoryImpl) AddCreatedBlock(block blockchain.Block) {
	bpr.Blocks = append(bpr.Blocks, block)
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
	for _, block := range bpr.Blocks {
		if string(block.GetSeal()) == blockId {
			return block, nil
		}
	}
	return nil, ErrNoStagedBlock
}

type CommitedBlockRepositoryImpl struct {
	mux *sync.RWMutex
	yggdrasill.BlockStorageManager
}

func (cbr *CommitedBlockRepositoryImpl) GetLastBlock() (blockchain.Block, error) {
	cbr.mux.Lock()
	defer cbr.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := cbr.BlockStorageManager.GetLastBlock(block)
	if err != nil {
		return nil, ErrGetCommitedBlock
	}

	return block, nil
}
func (cbr *CommitedBlockRepositoryImpl) GetBlockByHeight(height uint64) (blockchain.Block, error) {
	cbr.mux.Lock()
	defer cbr.mux.Unlock()

	block := &blockchain.DefaultBlock{}

	err := cbr.BlockStorageManager.GetBlockByHeight(block, height)
	if err != nil {
		return nil, ErrGetCommitedBlock
	}

	return block, nil
}
