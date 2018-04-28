package api

import (
	"github.com/it-chain/it-chain-Engine/txpool/domain/repository"
	"github.com/it-chain/it-chain-Engine/txpool/domain/service"
	"github.com/it-chain/it-chain-Engine/txpool/domain/model/transaction"
	"github.com/it-chain/it-chain-Engine/txpool/domain/model/timeout"
	"github.com/it-chain/it-chain-Engine/conf"
)

type TxpoolApi struct {
	txRepository  *repository.TransactionRepository
	timeoutTicker *timeout.TimeoutTicker
	maxTxByte     int
	messageApi    *service.MessageProducer
}

func NewTxpoolApi (txpoolRepo *repository.TransactionRepository, messageProducer *service.MessageProducer) *TxpoolApi{
	txpConfig := conf.GetConfiguration().Txpool

	return &TxpoolApi{
		txRepository:  txpoolRepo,
		timeoutTicker: timeout.NewTimeoutTicker(txpConfig.TimeoutMs),
		maxTxByte:     txpConfig.MaxTransactionByte,
		messageApi:    messageProducer,
	}
}

func (txpoolApi TxpoolApi) SaveTransaction(transaction transaction.Transaction) error {
	return (*txpoolApi.txRepository).Save(transaction)
}

func (txpoolApi TxpoolApi) RemoveTransaction(transactionId transaction.TransactionId) error {
	return (*txpoolApi.txRepository).Remove(transactionId)
}