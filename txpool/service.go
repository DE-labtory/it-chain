package txpool

import (
	"errors"
)

type TxpoolQueryService interface {
	GetLeader() Leader
	GetAllTransactions() ([]Transaction, error)
}

type TransferService interface {
	SendTransactionsToLeader(transactions []Transaction, leader Leader) error
}

type BlockService interface {
	ProposeBlock(transactions []Transaction) error
}

type TxTransferService struct {
	txpoolQueryService TxpoolQueryService
	transferService    TransferService
}

func NewTxPeriodicTransferService(queryService TxpoolQueryService, transferService TransferService) *TxTransferService {

	return &TxTransferService{
		txpoolQueryService: queryService,
		transferService:    transferService,
	}
}

func (t TxTransferService) TransferCreatedTxToLeader() error {

	//todo timeStamp check
	transactions, err := t.txpoolQueryService.GetAllTransactions()

	if err != nil {
		return err
	}

	leader := t.txpoolQueryService.GetLeader()

	if leader.LeaderId.ToString() == "" {
		return errors.New("there is no leader")
	}

	err = t.transferService.SendTransactionsToLeader(transactions, leader)

	return nil
}
