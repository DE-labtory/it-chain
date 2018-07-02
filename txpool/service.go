package txpool

type GrpcCommandService interface {
	SendLeaderTransactions(transactions []*Transaction, leader Leader) error
}

type BlockService interface {
	ProposeBlock(transactions []Transaction) error
}
