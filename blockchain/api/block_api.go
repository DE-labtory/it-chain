package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"fmt"
	"github.com/pkg/errors"
	"github.com/it-chain/midgard"
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
	eventRepository midgard.EventRepository
	publisherId          string
}

func NewBlockApi(blockRepository BlockRepository, eventRepository midgard.EventRepository, publisherId string) (BlockApi, error) {
	return BlockApi{
		blockRepository: blockRepository,
		eventRepository: eventRepository,
		publisherId:          publisherId,
	}, nil
}

// TODO: Check 과정에서 임의의 노드에게서 받은 blockchain 정보로 동기화 되었는지 확인한다.
func (bApi *BlockApi) SyncedCheck(block blockchain.Block) error {
	return nil
}
// 받은 block을 block pool에 추가한다.
func (bApi *BlockApi) AddBlockToPool(block blockchain.Block) error {
	if block == nil {
		fmt.Println("block is nil")
		return nil
	}

	pool := bApi.loadBlockPool()
	err := pool.Add(block)
	if err != nil {
		return err
	}
	return nil
}

// TODO
func (bApi *BlockApi) CheckAndSaveBlockFromPool(height blockchain.BlockHeight) error {
	pool := bApi.loadBlockPool()

	// Get block from pool
	blockFromPool := pool.Get(height)

	if blockFromPool == nil {
		return ErrNilBlock
	}

	// Get my last block
	lastBlock := &blockchain.DefaultBlock{}
	bApi.blockRepository.GetLastBlock(lastBlock)

	// Compare height
	if blockFromPool.GetHeight() > lastBlock.GetHeight() + 1 {
		// TODO: Start synchronize

	} else if blockFromPool.GetHeight() == lastBlock.GetHeight() + 1 {
		//
		bApi.blockRepository.AddBlock(blockFromPool)
		pool.Delete(blockFromPool)

	} else {
		// Got shorter height block, but this is not an error
		fmt.Printf("got shorter height block [%d < %d]", blockFromPool.GetHeight(), lastBlock.GetHeight());
	}

	return nil
}

func (bApi *BlockApi) loadBlockPool() blockchain.BlockPool {
	pool := blockchain.NewBlockPool()
	bApi.eventRepository.Load(pool, blockchain.BLOCK_POOL_AID)
	return pool
}
