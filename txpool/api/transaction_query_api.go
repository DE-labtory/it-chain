package api

import "github.com/it-chain/it-chain-Engine/txpool"

type TransactionQueryApi struct {
	txpool.TransactionRepository
}

func (t TransactionQueryApi) GetAllTransactions() ([]txpool.Transaction, error) {

	return nil, nil
}
