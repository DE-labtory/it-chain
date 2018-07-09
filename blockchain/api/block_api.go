package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/pkg/errors"
	"github.com/it-chain/midgard"
)

var ErrNilBlock = errors.New("block is nil")
var ErrSyncProcessing = errors.New("Sync is in progress")

type BlockApi struct {
	blockQueryApi blockchain.BlockQueryApi
	eventRepository midgard.EventRepository
	publisherId          string
}

func NewBlockApi(blockQueryApi blockchain.BlockQueryApi, eventRepository midgard.EventRepository, publisherId string) (BlockApi, error) {
	return BlockApi{
		blockQueryApi: blockQueryApi,
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
		return ErrNilBlock
	}

	pool := bApi.loadBlockPool()
	err := pool.Add(block)
	if err != nil {
		return err
	}
	return nil
}

func (bApi *BlockApi) CheckAndSaveBlockFromPool(height blockchain.BlockHeight) error {
	if bApi.SyncIsProgressing() {
		return ErrSyncProcessing
	}

	pool := bApi.loadBlockPool()
	// Get block from pool
	blockFromPool := pool.Get(height)
	if blockFromPool == nil {
		return ErrNilBlock
	}

	// Get my last block
	lastBlock := &blockchain.DefaultBlock{}
	bApi.blockQueryApi.GetLastBlock(lastBlock)

	// Compare height
	checkResult := getCheckResult(blockFromPool.GetHeight(), lastBlock.GetHeight())

	// Create action based on check result
	actionAfterCheck := blockchain.CreateActionAfterCheck(checkResult, pool, bApi.blockQueryApi)

	actionAfterCheck.DoAction(blockFromPool)

	return nil
}

func (bApi *BlockApi) SyncIsProgressing() blockchain.ProgressState {
	syncState := blockchain.NewBlockSyncState()
	bApi.eventRepository.Load(syncState, blockchain.BC_SYNC_STATE_AID)
	return syncState.IsProgressing()
}

func (bApi *BlockApi) loadBlockPool() blockchain.BlockPool {
	pool := blockchain.NewBlockPool()
	bApi.eventRepository.Load(pool, blockchain.BLOCK_POOL_AID)
	return pool
}

func getCheckResult(height1 uint64, height2 uint64) int64 {
	return int64(height1 - height2) - 1
}





