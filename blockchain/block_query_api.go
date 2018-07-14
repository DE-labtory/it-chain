package blockchain

// interface of api gateway query api
type BlockQueryApi interface {
	GetStagedBlockByHeight(blockHeight uint64) (Block, error)
	GetCommitedBlockByHeight(blockHeight uint64) (Block, error)
	GetStagedBlockById(blockId string) (Block, error)
	GetBlockByHeight(height BlockHeight) (Block, error)
	GetLastBlock() (Block, error)
	GetLastCommitedBlock() (Block, error)
}
