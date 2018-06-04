package txpool

type MessageDispatcher interface {
	SendLeaderTransactions(transactions []*Transaction, leader Leader) error
	ProposeBlock(transactions []Transaction) error
}
