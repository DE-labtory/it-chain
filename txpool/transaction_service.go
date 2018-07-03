package txpool

type TransactionService interface {
	TransferTxToLeader() error
}
