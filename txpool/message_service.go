package txpool

type MessageService interface {
	SendLeaderTransactions(transactions []*Transaction, leader Leader) error
}
