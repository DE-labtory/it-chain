package api_gateway

import "github.com/it-chain/it-chain-Engine/txpool"

// this is an api only for querying current state which is repository of transaction
type TransactionQueryApi struct {
	transactionRepository TransactionPoolRepository
}

// find all transactions that are created by not committed as a block
func (t TransactionQueryApi) FindUncommittedTransactions() []txpool.Transaction {
	return t.transactionRepository.FindAllTransaction()
}

// this repository is a current state of all uncommitted transactions
type TransactionPoolRepository interface {
	FindAllTransaction() []txpool.Transaction
	Save(transaction txpool.Transaction)
}

// this is an event_handler which listen all events related to transaction and update repository
// this struct will be relocated to other pkg
type TransactionEventListener struct {
	transactionRepository TransactionPoolRepository
}

// this function listens to TxCreatedEvent and update repository
func (t TransactionEventListener) HandleTransactionCreatedEvent(event txpool.TxCreatedEvent) {

	tx := event.GetTransaction()
	t.transactionRepository.Save(tx)
}
