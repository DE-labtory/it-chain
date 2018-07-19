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

type BlockCreateService interface {
	CreateBlock(block Block) error
}

type BlockCreateServiceImpl struct{}

func NewBlockCreateService() *BlockCreateServiceImpl {
	return &BlockCreateServiceImpl{}
}

func (bcs *BlockCreateServiceImpl) CreateBlock(block Block) error {
	// create BlockCreatedEvent
	// save it to event store
	return nil
}

type BlockStageService interface {
	StageBlock(block Block) error
}

type BlockStageServiceImpl struct{}

func NewBlockStageService() *BlockStageServiceImpl {
	return &BlockStageServiceImpl{}
}

func (bss *BlockStageServiceImpl) StageBlock(block Block) error {
	// create BlockStageEvent
	// save it to event store
	return nil
}

type BlockCommitService interface {
	CommitBlock(block Block) error
}

type BlockCommitServiceImpl struct{}

func NewBlockCommitService() *BlockCommitServiceImpl {
	return &BlockCommitServiceImpl{}
}

func (bcs *BlockCommitServiceImpl) CommitBlock(block Block) error {
	// create BlockCommitEvent
	// save it to event store
	return nil
}
