package txpool

//Transaction Repository interface
type TransactionRepository interface {
	Save(transaction Transaction) error
	Remove(id TransactionId) error
	FindById(id TransactionId) (Transaction, error)
	FindAll() ([]Transaction, error)
}

type LeaderRepository interface {
	GetLeader() Leader
	SetLeader(leader Leader)
}
