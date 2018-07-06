package infra

import "github.com/it-chain/it-chain-Engine/txpool"

//Transaction Repository interface
type TransactionRepository interface {
	Save(transaction txpool.Transaction) error
	Remove(id txpool.TransactionId) error
	FindById(id txpool.TransactionId) (txpool.Transaction, error)
	FindAll() ([]txpool.Transaction, error)
}
