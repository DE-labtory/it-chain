package txpool

type TxpoolQueryService interface {
	GetAllTransactions() ([]Transaction, error)
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

func NewBlockProposalService(queryService TxpoolQueryService, blockService BlockService, publisher Publisher) *BlockProposalService {

	return &BlockProposalService{
		publisher:          publisher,
		txpoolQueryService: queryService,
		blockService:       blockService,
	}
}

func (b BlockProposalService) ProposeBlock() error {

	// todo transaction size, number of tx
	transactions, err := b.txpoolQueryService.GetAllTransactions()

	if err != nil {
		return err
	}

	// todo filter txs based on time stamp
	//filterdTxs := filter(transactions, func(transaction Transaction) bool {
	//
	//	return
	//})

	return b.blockService.ProposeBlock(transactions)
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
