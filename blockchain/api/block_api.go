package api

import (
	"github.com/it-chain/engine/blockchain"
)

type BlockApi struct {
	publisherId string
}

func NewBlockApi(publisherId string) (BlockApi, error) {
	return BlockApi{
		publisherId: publisherId,
	}, nil
}

// TODO: Check 과정에서 임의의 노드에게서 받은 blockchain 정보로 동기화 되었는지 확인한다.
func (bApi *BlockApi) SyncedCheck(block blockchain.Block) error {
	return nil
}

// 받은 block을 block pool에 추가한다.
func (bApi *BlockApi) AddBlockToPool(block blockchain.Block) error {
	return nil
}

func (bApi *BlockApi) CheckAndSaveBlockFromPool(height blockchain.BlockHeight) error {
	return nil
}

func (bApi *BlockApi) SyncIsProgressing() blockchain.ProgressState {
	return blockchain.DONE
}

func (bApi *BlockApi) loadBlockPool() blockchain.BlockPool {
	return nil
}

func compareHeight(height1 uint64, height2 uint64) int64 {
	return int64(height1-height2) - 1
}
