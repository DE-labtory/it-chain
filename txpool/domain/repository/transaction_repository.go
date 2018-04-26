package repository

import (
	"github.com/it-chain/it-chain-Engine/txpool/domain/model/transaction"
)

type TransactionRepository interface {
	Save(transaction transaction.Transaction) error
	Remove(id transaction.TransactionId) error
	FindById(id transaction.TransactionId) (*transaction.Transaction, error)
	FindAll() ([]*transaction.Transaction, error)
}
