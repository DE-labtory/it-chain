package blockchaindb

import (
	"it-chain/service/blockchain"
)

// HistoryDB - an interface that a history database should implement
// todo 다른 기능들 추가
type BlockChainDB interface {
	Commit(block *blockchain.Block) error
}