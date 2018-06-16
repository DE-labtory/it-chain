package txpool

type BlockService interface {
	ProposeBlock(transactions []Transaction) error
}
