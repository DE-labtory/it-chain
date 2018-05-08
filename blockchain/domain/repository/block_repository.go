package repository

import (
	"github.com/it-chain/it-chain-Engine/blockchain/domain/model/block"
	"github.com/it-chain/it-chain-Engine/blockchain/domain/model/transaction"
)

type BlockRepository interface {
	Close()
	AddBlock(block block.Block) error
	GetBlockByNumber(blockNumber uint64) (block.Block, error)
	GetBlockByHash(hash string) (block.Block, error)
	GetLastBlock() (block.Block, error)
	GetTransactionByTxID(txid string) (transaction.Trasaction, error)
	GetBlockByTxID(txid string) (block.Block, error)
}
