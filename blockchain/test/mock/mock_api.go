package mock

import "github.com/it-chain/it-chain-Engine/blockchain"

type BlockApi struct {
	AddBlockToPoolFunc            func(block blockchain.Block) error
	CommitBlockFromPoolOrSyncFunc func(blockId string) error
}

func (api BlockApi) AddBlockToPool(block blockchain.Block) error {
	return api.AddBlockToPoolFunc(block)
}

func (api BlockApi) CommitBlockFromPoolOrSync(blockId string) error {
	return api.CommitBlockFromPoolOrSyncFunc(blockId)
}

type MockSyncBlockApi struct {
	SyncedCheckFunc func(block blockchain.Block) error
}

func (ba MockSyncBlockApi) SyncedCheck(block blockchain.Block) error {
	return ba.SyncedCheckFunc(block)
}
