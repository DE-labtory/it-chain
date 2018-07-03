package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"fmt"
	"github.com/pkg/errors"
)

type BlockRepository interface {
	GetValidator() blockchain.Validator
	GetLastBlock(block blockchain.Block) error
	AddBlock(block blockchain.Block) error
	NewEmptyBlock() (blockchain.Block, error)
	Close()
}

var ErrNilBlock = errors.New("block is nil")

type BlockApi struct {
	blockRepository BlockRepository
	publisherId          string
	blockPool blockchain.BlockPool
}

func NewBlockApi(blockRepository BlockRepository, publisherId string, blockPool blockchain.BlockPool) (BlockApi, error) {
	return BlockApi{
		blockRepository: blockRepository,
		publisherId:          publisherId,
		blockPool: blockPool,
	}, nil
}

// TODO: Check 과정에서 임의의 노드에게서 받은 blockchain 정보로 동기화 되었는지 확인한다.
func (bApi *BlockApi) SyncedCheck(block blockchain.Block) error {
	return nil
}
// 받은 block을 block pool에 추가한다.
func (bApi *BlockApi) AddBlockToPool(block blockchain.Block) {
	if block == nil {
		fmt.Println("block is nil")
		return
	}
	bApi.blockPool.Add(block)
}

// TODO
func (bApi *BlockApi) CheckAndSaveBlockFromPool(height blockchain.BlockHeight) error {
	// Get block from pool
	block := bApi.blockPool.Get(height)
	if block == nil {
		return ErrNilBlock
	}

	// Get my last block
	lastBlock := &blockchain.DefaultBlock{}
	bApi.blockRepository.GetLastBlock(lastBlock)

	// Compare height
	if block.GetHeight() > lastBlock.GetHeight() + 1 {
		// TODO: Start synchronize

	} else if block.GetHeight() == lastBlock.GetHeight() + 1 {
		// Save
		bApi.blockRepository.AddBlock(block)

		bApi.blockPool.Delete(height)

	} else {
		// Got shorter height block, but this is not an error
		fmt.Printf("got shorter height block [%d < %d]", block.GetHeight(), lastBlock.GetHeight());
	}

	return nil
}
