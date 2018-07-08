package api_gateway

import "github.com/it-chain/it-chain-Engine/txpool"

type TransactionQueryApi struct {
	transactionRepository TransactionPoolRepository
}

func (t TransactionQueryApi) FindUncommittedTransactions() []txpool.Transaction {
	return t.transactionRepository.FindAllTransaction()
}

type TransactionPoolRepository interface {
	FindAllTransaction() []txpool.Transaction
	Save(transaction txpool.Transaction)
}

// this is an event_handler which listen all events related to transaction and update repository
type TransactionEventListener struct {
	transactionRepository TransactionPoolRepository
}

func (t TransactionEventListener) HandleTransactionCreatedEvent(event txpool.TxCreatedEvent) {

	tx := event.GetTransaction()
	t.transactionRepository.Save(tx)
}
