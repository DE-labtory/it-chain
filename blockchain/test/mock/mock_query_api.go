package mock

import "github.com/it-chain/it-chain-Engine/blockchain"

type BlockQueryApi struct {
	GetLastBlockFunc func() (blockchain.Block, error)
	GetBlockByHeightFunc func(blockHeight uint64) (blockchain.Block, error)
	GetBlockBySealFunc func(seal []byte) (blockchain.Block, error)
	GetBlockByTxIDFunc func(txid string) (blockchain.Block, error)
	GetTransactionByTxIDFunc func(txid string) (blockchain.Transaction, error)
}
func (br BlockQueryApi) GetLastBlock() (blockchain.Block, error) {
	return br.GetLastBlockFunc()
}
func (br BlockQueryApi) GetBlockByHeight(blockHeight uint64) (blockchain.Block, error) {
	return br.GetBlockByHeightFunc(blockHeight)
}
func (br BlockQueryApi) GetBlockBySeal(seal []byte) (blockchain.Block, error) {
	return br.GetBlockBySealFunc(seal)
}
func (br BlockQueryApi) GetBlockByTxID(txid string) (blockchain.Block, error) {
	return br.GetBlockByTxIDFunc(txid)
}
func (br BlockQueryApi) GetTransactionByTxID(txid string) (blockchain.Transaction, error) {
	return br.GetTransactionByTxIDFunc(txid)
}


