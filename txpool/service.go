package txpool

type MessageService interface {
	SendLeaderTransactions(transactions []*Transaction, leader Leader) error
}

type BlockService interface {
	ProposeBlock(transactions []Transaction) error
}
