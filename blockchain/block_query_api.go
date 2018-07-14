package blockchain

type BlockQueryApi interface {
	GetStagedBlockByHeight(height uint64) (Block, error)
	GetStagedBlockById(blockId string) (Block, error)
	GetLastCommitedBlock() (Block, error)
	GetCommitedBlockByHeight(height uint64) (Block, error)
}
