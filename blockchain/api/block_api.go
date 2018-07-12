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
	blockQueryApi   blockchain.BlockQueryApi
	eventRepository midgard.EventRepository
	publisherId     string
	SyncService     blockchain.SyncService
}

func NewBlockApi(blockQueryApi blockchain.BlockQueryApi, eventRepository midgard.EventRepository, publisherId string) (BlockApi, error) {
	return BlockApi{
		blockQueryApi:   blockQueryApi,
		eventRepository: eventRepository,
		publisherId:     publisherId,
	}, nil
}

func (bApi *BlockApi) Synchronize() {}

// TODO: Check 과정에서 임의의 노드에게서 받은 blockchain 정보로 동기화 되었는지 확인한다.
func (bApi *BlockApi) SyncedCheck(block blockchain.Block) error {
	return nil
}

func (bApi *BlockApi) construct() error {

	// 0. 싱크 프로세싱 스테이트 변경: True SetProgress(Processing)
	syncState := blockchain.NewBlockSyncState()
	syncState.SetProgress(blockchain.PROGRESSING)

	// 0.5 Peer객체 셋팅: 다른 함수 써야함.
	anonymousPeerId := blockchain.PeerId{Id: "anonymous"}

	// 1. cashe에 상대방의 lastBlock의 height 저장.
	yourLastBlock, err := bApi.blockQueryApi.GetLastBlock()
	if err != nil {
		return err
	}
	// 2. 나의 라스트 헤이트 세팅: 이벤트 스토어에서 나의 lastBlock 가져옴. 하나도 없을 경우 0, 있다면 마지막.
	myLastBlock, err := bApi.blockQueryApi.GetLastBlock()
	if err != nil {
		return err
	}

	myLastHeight := myLastBlock.GetHeight()

	// 3. while 문 사용: 나의 lastHeight와 상대방의 last height 비교
	for yourLastBlock.GetHeight() <= myLastHeight {

		// 4. 리퀘스트 블록: 내가 하나도 없을 경우에 0, 있다면 lastHeight + 1을 요청
		bApi.SyncService.RequestBlock(anonymousPeerId, myLastHeight)
		// 5. 이벤트 스토어에 세이브.
		blockchain.CommitBlock(nil)

		// 6. 나의 last height 다시 세팅.
		myLastHeight++

	}

	// 7. 싱크 프로세싱 스테이트 변경: False
	syncState.SetProgress(blockchain.DONE)

	return nil
}

func (bApi *BlockApi) SaveBlock(block blockchain.Block) error {
	if block == nil {
		return ErrNilBlock
	}

	action := blockchain.NewSaveAction()
	action.DoAction(block)
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
