package txpool

import (
	"errors"
	"log"
)

type GrpcCommandService interface {
	SendLeaderTransactions(transactions []*Transaction, leader Leader) error
}

type BlockService interface {
	ProposeBlock(transactions []Transaction) error
}

type TxPeriodicTransferService struct {
	txRepository       TransactionRepository
	leaderRepository   LeaderRepository
	grpcCommandService GrpcCommandService
}

func NewTxPeriodicTransferService(tr TransactionRepository, lr LeaderRepository, gcs GrpcCommandService) *TxPeriodicTransferService {
	return &TxPeriodicTransferService{
		txRepository:       tr,
		leaderRepository:   lr,
		grpcCommandService: gcs,
	}
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

	err = t.grpcCommandService.SendLeaderTransactions(transactions, leader)

	if err != nil {
		log.Println(err.Error())
		return err
	}

	//if err := t.removeTxs(transactions); err != nil {
	//	log.Println(err.Error())
	//	return err
	//}

	return nil
}

//func (t TxPeriodicTransferService) removeTxs(transactions []*Transaction) error {
//
//	for _, tx := range transactions {
//		if err := t.txRepository.Remove(tx.TxId); err != nil {
//			return err
//		}
//	}
//
//	return nil
//}
