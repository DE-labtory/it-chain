package blockchaindb

import (
	"it-chain/service/blockchain"
)


// HistoryDB - an interface that a history database should implement
type BlockChainDB interface {
	Commit(block *blockchain.Block) error
}