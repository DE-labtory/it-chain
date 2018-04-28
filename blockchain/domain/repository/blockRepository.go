package repository

import (
	"github.com/it-chain/it-chain-Engine/blockchain/domain/model/block"
	"github.com/it-chain/yggdrasill/transaction"
)

type BlockRepository interface {
	Close()
	AddBlock(block block.Block) error
	GetLastBlock() block.Block
	GetTransactionsById(id string) transaction.Transaction
}
