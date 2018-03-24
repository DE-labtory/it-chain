package db

import "github.com/it-chain/it-chain-Engine/blockchain/domain/model"

type BlockChainDB interface {
	Close()
	AddBlock(block *model.Block) error
	AddUnconfirmedBlock(block *model.Block) error
	GetBlockByNumber(blockNumber uint64) (*model.Block, error)
	GetBlockByHash(hash string) (*model.Block, error)
	GetLastBlock() (*model.Block, error)
	GetTransactionByTxID(txid string) (*model.Transaction, error)
	GetBlockByTxID(txid string) (*model.Block, error)
}