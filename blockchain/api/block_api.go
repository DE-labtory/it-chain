package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/midgard"
)

var ErrNilBlock = errors.New("block is nil")
var ErrSyncProcessing = errors.New("Sync is in progress")
var ErrGetLastBlock = errors.New("failed get last block")

type BlockApi struct {
	syncService     blockchain.SyncService
	peerService     blockchain.PeerService
	eventRepository midgard.EventRepository
	publisherId     string
}

func NewBlockApi(syncService blockchain.SyncService, peerService blockchain.PeerService, eventRepository midgard.EventRepository, publisherId string) (BlockApi, error) {
	return BlockApi{
		syncService:     syncService,
		peerService:     peerService,
		eventRepository: eventRepository,
		publisherId:     publisherId,
	}, nil
}

// ToDo: 임의의 노드와 블록 체인을 동기화합니다.
func (bApi *BlockApi) Synchronize() error {

	// 싱크 프로세싱 스테이트 변경
	syncState := blockchain.NewBlockSyncState()
	syncState.SetProgress(blockchain.PROGRESSING)

	// peer random으로 가져오기
	peer, err := bApi.peerService.GetRandomPeer()

	if err != nil {
		return err
	}

	// sync 하기
	if err = bApi.syncService.Sync(peer); err != nil {
		return err
	}

	syncState.SetProgress(blockchain.DONE)

	//event.save

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

func (bApi *BlockApi) retrieveAndCommitBlock() {

}

func (bApi *BlockApi) CheckAndSaveBlockFromPool(height blockchain.BlockHeight) error {
	if bApi.SyncIsProgressing() {
		return ErrSyncProcessing
	}

	pool := bApi.loadBlockPool()

	blockFromPool := pool.Get(height)
	if blockFromPool == nil {
		return ErrNilBlock
	}

	lastBlock, err := bApi.blockQueryApi.GetLastBlock()
	if err != nil {
		return ErrGetLastBlock
	}

	isSafeToSave := compareHeight(blockFromPool.GetHeight(), lastBlock.GetHeight())

	action := blockchain.CreateSaveOrSyncAction(isSafeToSave)

	action.DoAction(blockFromPool)

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

func compareHeight(height1 uint64, height2 uint64) int64 {
	return int64(height1-height2) - 1
}
