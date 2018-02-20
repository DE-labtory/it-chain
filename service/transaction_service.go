package service

import (
	"it-chain/domain"
)

type TransactionService interface {
	AddTransaction(tx *domain.Transaction) error
	DeleteTransactions(txs []*domain.Transaction) error
	GetTransactions(limit int) ([]*domain.Transaction, error)
	CreateTransaction(txData *domain.TxData) (*domain.Transaction, error)
}