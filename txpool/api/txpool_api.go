package api

import (
	"time"

	"github.com/it-chain/it-chain-Engine/txpool/domain/repository"
	"github.com/it-chain/it-chain-Engine/txpool/domain/service"
	"github.com/it-chain/it-chain-Engine/txpool/domain/model/transaction"
	"github.com/it-chain/it-chain-Engine/conf"
)

// todo api 만들어라 준희야
type TxpoolApi struct {
	txRepository  repository.TransactionRepository
	timeoutTicker *time.Ticker
	maxTxNum      int
	messageApi    service.MessageProducer
}

func NewTxpoolApi () *TxpoolApi{
	txpoolConfig = conf.GetConfiguration().Txpool

	return &TxpoolApi{
		txRepository:  nil,
		timeoutTicker: nil,
		maxTxNum:      0,
		messageApi:    nil,
	}
}

func (txpoolApi TxpoolApi) SaveTransaction(transaction transaction.Transaction) error {
	return txpoolApi.txRepository.Save(transaction)
}
