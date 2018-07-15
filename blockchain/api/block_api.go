package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type BlockApi struct {
	syncService       blockchain.SyncService
	peerService       blockchain.PeerService
	blockQueryService blockchain.BlockQueryService
	eventRepository   midgard.EventRepository
	publisherId       string
}

func NewBlockApi(syncService blockchain.SyncService, peerService blockchain.PeerService, blockQueryService blockchain.BlockQueryService, eventRepository midgard.EventRepository, publisherId string) (BlockApi, error) {
	return BlockApi{
		syncService:       syncService,
		peerService:       peerService,
		blockQueryService: blockQueryService,
		eventRepository:   eventRepository,
		publisherId:       publisherId,
	}, nil
}

// Synchronize blockchain with a peer in p2p network
func (bApi *BlockApi) Synchronize() error {

	// Set syncState : Progressing
	syncState := blockchain.NewBlockSyncState()
	syncState.SetProgress(blockchain.PROGRESSING)

	// Get a random peer in p2p network
	peer, err := bApi.peerService.GetRandomPeer()

	if err != nil {
		return ErrGetRandomPeer
	}

	// Synchronize blockchain with a random peer
	if err = bApi.syncService.SyncWithPeer(peer); err != nil {
		return ErrSyncWithPeer
	}

	// Set syncState : Done
	syncState.SetProgress(blockchain.DONE)

	return nil
}

// 받은 block을 block pool에 추가한다.
func (bApi *BlockApi) StageBlock(block blockchain.Block) error {
	if block == nil {
		return ErrNilBlock
	}

	err := blockchain.StageBlock(block)
	if err != nil {
		return err
	}

	return nil
}

func (bApi *BlockApi) CommitBlockFromPoolOrSync(blockId string) error {
	if bApi.SyncIsProgressing() {
		return ErrSyncProcessing
	}

	block, err := bApi.blockQueryService.GetStagedBlockById(blockId)
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
	eventstore.Load(syncState, blockchain.BC_SYNC_STATE_AID)
	return syncState.IsProgressing()
}

func (bApi *BlockApi) CompareLastBlockHeightWith(targetBlock blockchain.Block) (uint64, error) {
	lastBlock, err := bApi.blockQueryService.GetLastCommitedBlock()
	if err != nil {
		return 0, ErrGetLastBlock
	}
	return targetBlock.GetHeight() - lastBlock.GetHeight(), nil
}
