package api

import (
	"time"

	"github.com/it-chain/it-chain-Engine/txpool/domain/repository"
	"github.com/it-chain/it-chain-Engine/txpool/domain/service"
)

// todo api 만들어라 준희야
type TxpoolApi struct {
	txRepository  repository.TransactionRepository
	timeoutTicker *time.Ticker
	maxTxNum      int
	messageApi    service.MessageProducer
}
