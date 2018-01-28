package blockchaindb

import (
	"it-chain/service/blockchain"
)

type BlockChainDB interface {
	Close()
	AddBlock(block *blockchain.Block) error
	AddUnconfirmedBlock(block *blockchain.Block) error
	GetBlockByNumber(blockNumber uint64) (*blockchain.Block, error)
	GetBlockByHash(hash string) (*blockchain.Block, error)
	GetLastBlock() (*blockchain.Block, error)
	GetTransactionByTxID(txid string) (*blockchain.Transaction, error)
	GetBlockByTxID(txid string) (*blockchain.Block, error)
}