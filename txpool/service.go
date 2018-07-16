package txpool

import "log"

type TxpoolQueryService interface {
	FindUncommittedTransactions() ([]Transaction, error)
}

type TransferService interface {
	SendTransactionsToLeader(transactions []Transaction, leader Leader) error
}

type BlockService interface {
	ProposeBlock(transactions []Transaction) error
}

type BlockProposalService struct {
	txpoolQueryService TxpoolQueryService
	blockService       BlockService
	publisher          Publisher
}

type Publisher func(exchange string, topic string, data interface{}) (err error)

func NewBlockProposalService(queryService TxpoolQueryService, blockService BlockService) *BlockProposalService {

	return &BlockProposalService{
		txpoolQueryService: queryService,
		blockService:       blockService,
	}
}

// todo do not delete transaction immediately
// todo transaction will be deleted when block are committed
func (b BlockProposalService) ProposeBlock() error {

	// todo transaction size, number of tx
	transactions, err := b.txpoolQueryService.FindUncommittedTransactions()

	log.Printf("proposing transactions [%s]", transactions)

	if err != nil {
		return err
	}

	// todo filter txs based on time stamp
	//filterdTxs := filter(transactions, func(transaction Transaction) bool {
	//
	//	return
	//})

	if len(transactions) == 0 {
		return nil
	}

	err = b.blockService.ProposeBlock(transactions)

	if err != nil {
		return err
	}

	for _, tx := range transactions {
		DeleteTransaction(tx)
	}

	return nil
}

func filter(vs []Transaction, f func(Transaction) bool) []Transaction {
	vsf := make([]Transaction, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}
