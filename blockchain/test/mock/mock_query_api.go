package mock

import "github.com/it-chain/it-chain-Engine/blockchain"

type MockBlockQueryApi struct {
	GetLastBlockFunc func() (blockchain.Block, error)
	GetBlockByHeightFunc func(blockHeight uint64) (blockchain.Block, error)
	GetBlockBySealFunc func(seal []byte) (blockchain.Block, error)
	GetBlockByTxIDFunc func(txid string) (blockchain.Block, error)
	GetTransactionByTxIDFunc func(txid string) (blockchain.Transaction, error)
}
func (br MockBlockQueryApi) GetLastBlock() (blockchain.Block, error) {
	return br.GetLastBlockFunc()
}
func (br MockBlockQueryApi) GetBlockByHeight(blockHeight uint64) (blockchain.Block, error) {
	return br.GetBlockByHeightFunc(blockHeight)
}
func (br MockBlockQueryApi) GetBlockBySeal(seal []byte) (blockchain.Block, error) {
	return br.GetBlockBySealFunc(seal)
}
func (br MockBlockQueryApi) GetBlockByTxID(txid string) (blockchain.Block, error) {
	return br.GetBlockByTxIDFunc(txid)
}
func (br MockBlockQueryApi) GetTransactionByTxID(txid string) (blockchain.Transaction, error) {
	return br.GetTransactionByTxIDFunc(txid)
}


