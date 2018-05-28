package blockchain

type BlockRepository interface {
	Close()
	AddBlock(block Block) error
	GetBlockByNumber(blockNumber uint64) (Block, error)
	GetBlockByHash(hash string) (Block, error)
	GetLastBlock() (Block, error)
	GetTransactionByTxID(txid string) (Transaction, error)
	GetBlockByTxID(txid string) (Block, error)
}
