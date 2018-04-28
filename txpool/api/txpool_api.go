package api

import (
	"github.com/it-chain/it-chain-Engine/txpool/domain/repository"
	"github.com/it-chain/it-chain-Engine/txpool/domain/service"
	"github.com/it-chain/it-chain-Engine/txpool/domain/model/transaction"
	"github.com/it-chain/it-chain-Engine/txpool/domain/model/timeout"
	"github.com/it-chain/it-chain-Engine/conf"
)

// todo api 만들어라 준희야
type TxpoolApi struct {
	txRepository  repository.TransactionRepository
	timeoutTicker *timeout.TimeoutTicker
	maxTxByte     int
	messageApi    service.MessageProducer
}

// TODO assign txRepo, msgApi
func NewTxpoolApi () *TxpoolApi{
	txpConfig := conf.GetConfiguration().Txpool

	return &TxpoolApi{
		txRepository:  nil,
		timeoutTicker: timeout.NewTimeoutTicker(txpConfig.TimoutMs),
		maxTxByte:     txpConfig.MaxTransactionByte,
		messageApi:    nil,
	}
}

// TODO impl
func (txpoolApi TxpoolApi) SaveTransaction(transaction transaction.Transaction) error {
	return txpoolApi.txRepository.Save(transaction)
}
