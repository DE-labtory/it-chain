package mock

import "github.com/it-chain/it-chain-Engine/blockchain"

type MockBlockApi struct {
	AddBlockToPoolFunc func(block blockchain.Block) error
	CheckAndSaveBlockFromPoolFunc func(height blockchain.BlockHeight) error
}
func (api MockBlockApi) AddBlockToPool(block blockchain.Block) error {
	return api.AddBlockToPoolFunc(block)
}

func (api MockBlockApi) CheckAndSaveBlockFromPool(height blockchain.BlockHeight) error {
	return api.CheckAndSaveBlockFromPoolFunc(height)
}

type MockSyncBlockApi struct {
	SyncedCheckFunc func(block blockchain.Block) error
}

func (ba MockSyncBlockApi) SyncedCheck(block blockchain.Block) error {
	return ba.SyncedCheckFunc(block)
}