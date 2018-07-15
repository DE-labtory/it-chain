package mock

import "github.com/it-chain/it-chain-Engine/blockchain"

type BlockQueryApi struct {
	GetLastBlockFunc             func() (blockchain.Block, error)
	GetBlockByHeightFunc         func(blockHeight uint64) (blockchain.Block, error)
	GetStagedBlockByHeightFunc   func(blockHeight uint64) (blockchain.Block, error)
	GetStagedBlockByIdFunc       func(blockId string) (blockchain.Block, error)
	GetLastCommitedBlockFunc     func() (blockchain.Block, error)
	GetCommitedBlockByHeightFunc func(blockHeight uint64) (blockchain.Block, error)
}

func (br BlockQueryApi) GetLastBlock() (blockchain.Block, error) {
	return br.GetLastBlockFunc()
}
func (br BlockQueryApi) GetBlockByHeight(blockHeight uint64) (blockchain.Block, error) {
	return br.GetBlockByHeightFunc(blockHeight)
}
func (br BlockQueryApi) GetStagedBlockByHeight(blockHeight uint64) (blockchain.Block, error) {
	return br.GetStagedBlockByHeightFunc(blockHeight)
}
func (br BlockQueryApi) GetStagedBlockById(blockId string) (blockchain.Block, error) {
	return br.GetStagedBlockByIdFunc(blockId)
}
func (br BlockQueryApi) GetLastCommitedBlock() (blockchain.Block, error) {
	return br.GetLastCommitedBlockFunc()
}
func (br BlockQueryApi) GetCommitedBlockByHeight(blockHeight uint64) (blockchain.Block, error) {
	return br.GetCommitedBlockByHeightFunc(blockHeight)
}
