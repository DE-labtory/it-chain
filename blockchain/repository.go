package blockchain

type BlockRepository interface {
	Save(block DefaultBlock) error
	FindLast() (DefaultBlock, error)
	FindByHeight(height BlockHeight) (DefaultBlock, error)
	FindBySeal(seal string) (DefaultBlock, error)
	FindAll() ([]DefaultBlock, error)
}
