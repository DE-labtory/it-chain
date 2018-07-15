package blockchain

type BlockService interface {
	ExecuteBlock(block Block) error
}
