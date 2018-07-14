package api

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/core/eventstore"
)

type BlockApi struct {
	blockQueryApi blockchain.BlockQueryApi
	publisherId   string
}

func NewBlockApi(blockQueryApi blockchain.BlockQueryApi, publisherId string) (BlockApi, error) {
	return BlockApi{
		blockQueryApi: blockQueryApi,
		publisherId:   publisherId,
	}, nil
}

// TODO: blockchain이 동기화 되었는지 확인한다.
func (bApi *BlockApi) SyncedCheck(block blockchain.Block) error {
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

	block, err := bApi.blockQueryApi.GetStagedBlockById(blockId)
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
	lastBlock, err := bApi.blockQueryApi.GetLastCommitedBlock()
	if err != nil {
		return 0, ErrGetLastBlock
	}
	return targetBlock.GetHeight() - lastBlock.GetHeight(), nil
}
