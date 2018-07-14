package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/midgard"
	"github.com/pkg/errors"
)

var ErrNilBlock = errors.New("block is nil")
var ErrSyncProcessing = errors.New("Sync is in progress")
var ErrGetLastBlock = errors.New("failed get last block")

type BlockApi struct {
	blockQueryApi   blockchain.BlockQueryApi
	eventRepository midgard.EventRepository
	publisherId     string
}

func NewBlockApi(blockQueryApi blockchain.BlockQueryApi, eventRepository midgard.EventRepository, publisherId string) (BlockApi, error) {
	return BlockApi{
		blockQueryApi:   blockQueryApi,
		eventRepository: eventRepository,
		publisherId:     publisherId,
	}, nil
}

// TODO: blockchain이 동기화 되었는지 확인한다.
func (bApi *BlockApi) SyncedCheck(block blockchain.Block) error {
	return nil
}

// 받은 block을 block pool에 추가한다.
func (bApi *BlockApi) AddBlockToPool(block blockchain.Block) error {
	if block == nil {
		return ErrNilBlock
	}

	err := blockchain.SaveBlockStagedEventToEventStore(block)
	if err != nil {
		return err
	}

	return nil
}

func (bApi *BlockApi) CommitBlockFromPoolOrSync(height blockchain.BlockHeight) error {
	if bApi.SyncIsProgressing() {
		return ErrSyncProcessing
	}

	block, err := bApi.blockQueryApi.GetBlockByHeight(height)
	if err != nil {
		return err
	}

	diff, err := bApi.CompareLastBlockHeightWith(block)
	if err != nil {
		return err
	}

	bApi.commitBlockOrSyncByHeightDifference(diff, block)

	return nil
}

func (bApi *BlockApi) commitBlockOrSyncByHeightDifference(diff uint64, block blockchain.Block) error {
	if diff > 1 {
		// TODO: Start synchronize
	} else if diff == 1 {
		return blockchain.CommitBlock(block)
	}
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

func (bApi *BlockApi) CompareLastBlockHeightWith(targetBlock blockchain.Block) (uint64, error) {
	lastBlock, err := bApi.blockQueryApi.GetLastBlock()
	if err != nil {
		return 0, ErrGetLastBlock
	}

	return targetBlock.GetHeight() - lastBlock.GetHeight(), nil
}
