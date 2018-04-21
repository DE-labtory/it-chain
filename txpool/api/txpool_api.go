package api

import (
	"time"

	"github.com/it-chain/it-chain-Engine/txpool/domain/repository"
	"github.com/it-chain/it-chain-Engine/txpool/domain/service"
)

type TxpoolApi struct {
	txRepository  repository.TransactionRepository
	timeoutTicker *time.Ticker
	maxTxNum      int
	messageApi    service.MessageService
}
