package txpool

type MessageDispatcher interface {
	SendGrpcTransactions(transactions []*Transaction, leader Leader) error
	ProposeBlock(transactions []Transaction) error
}
