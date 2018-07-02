package batch

import (
"errors"
"log"
	"github.com/it-chain/it-chain-Engine/txpool"
	"time"
)

type TxPeriodicTransferService struct {
	txRepository     txpool.TransactionRepository
	leaderRepository txpool.LeaderRepository
	grpcCommandService   txpool.GrpcCommandService
	batcher *txpool.TimeoutBatcher
}

func NewTxPeriodicTransferService(tr txpool.TransactionRepository, lr txpool.LeaderRepository, gcs txpool.GrpcCommandService) *TxPeriodicTransferService {
	return &TxPeriodicTransferService{
		txRepository: tr,
		leaderRepository: lr,
		grpcCommandService: gcs,
		batcher: txpool.GetTimeOutBatcherInstance(),
	}
}

func (t TxPeriodicTransferService) Run(duration time.Duration) chan struct{}{
	return t.batcher.Run(t.TransferTxToLeader, duration)
}

func (t TxPeriodicTransferService) TransferTxToLeader() error {
	transactions, err := t.txRepository.FindAll()

	if err != nil {
		log.Println(err.Error())
		return err
	}

	leader := t.leaderRepository.GetLeader()

	if leader.StringLeaderId() == "" {
		return errors.New("there is no leader")
	}

	if err := t.removeTxs(transactions); err != nil {
		log.Println(err.Error())
		return err
	}

	err = t.grpcCommandService.SendLeaderTransactions(transactions, leader)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (t TxPeriodicTransferService) removeTxs(transactions []*txpool.Transaction) error {
	for _, tx := range transactions {
		if err := t.txRepository.Remove(tx.TxId); err != nil {
			return err
		}
	}
	return nil
}