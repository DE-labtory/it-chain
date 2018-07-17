package blockchain

type BlockQueryService interface {
	BlockQueryInnerService
}

type BlockQueryInnerService interface {
	GetStagedBlockByHeight(height BlockHeight) (Block, error)
	GetStagedBlockById(blockId string) (Block, error)
	GetLastCommitedBlock() (Block, error)
	GetCommitedBlockByHeight(height BlockHeight) (Block, error)
}
