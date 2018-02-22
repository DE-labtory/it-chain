package blockchaindb

import (
	"it-chain/domain"
)

type BlockChainDB interface {
	Close()
	AddBlock(block *domain.Block) error
	AddUnconfirmedBlock(block *domain.Block) error
	GetBlockByNumber(blockNumber uint64) (*domain.Block, error)
	GetBlockByHash(hash string) (*domain.Block, error)
	GetLastBlock() (*domain.Block, error)
	GetTransactionByTxID(txid string) (*domain.Transaction, error)
	GetBlockByTxID(txid string) (*domain.Block, error)
}