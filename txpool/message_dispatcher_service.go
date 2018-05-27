package txpool

type MessageDispatcher interface {
	SendTransactions(transactions []Transaction, leader Leader) error
	ProposeBlock(transactions []Transaction) error
}
