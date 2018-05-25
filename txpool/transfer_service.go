package txpool

import (
	"github.com/it-chain/midgard"
)

type TxPeriodicTransferService struct {
	txRepository    TransactionRepository
	eventRepository midgard.Repository
}

func (t TxPeriodicTransferService) TransferTx() {

}
