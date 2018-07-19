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
	engineMode         string
	txpoolQueryService TxpoolQueryService
	blockService       BlockService
}

type Publisher func(exchange string, topic string, data interface{}) (err error)

func NewBlockProposalService(queryService TxpoolQueryService, blockService BlockService, engineMode string) *BlockProposalService {

	return &BlockProposalService{
		txpoolQueryService: queryService,
		blockService:       blockService,
		engineMode:         engineMode,
	}
}

// todo do not delete transaction immediately
// todo transaction will be deleted when block are committed
func (b BlockProposalService) ProposeBlock() error {

	// todo transaction size, number of tx
	transactions, err := b.txpoolQueryService.FindUncommittedTransactions()

	if err != nil {
		return err
	}

	if len(transactions) == 0 {
		return nil
	}

	if b.engineMode == "solo" {
		//propose transaction when solo mode
		err = b.blockService.ProposeBlock(transactions)

		if err != nil {
			return err
		}

		log.Printf("transactions are proposed [%v]", transactions)

		for _, tx := range transactions {
			DeleteTransaction(tx)
		}

		return nil
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
