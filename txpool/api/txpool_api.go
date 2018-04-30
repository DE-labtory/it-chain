package api

import (
	"github.com/it-chain/it-chain-Engine/txpool/domain/repository"
	"github.com/it-chain/it-chain-Engine/txpool/domain/service"
	"github.com/it-chain/it-chain-Engine/txpool/domain/model/transaction"
	"github.com/it-chain/it-chain-Engine/txpool/domain/model/timeout"
	"github.com/it-chain/it-chain-Engine/conf"
	"errors"
	"github.com/it-chain/it-chain-Engine/common"
)

type TxpoolApi struct {
	txRepository  *repository.TransactionRepository
	timeoutTicker *timeout.TimeoutTicker
	maxTxByte     int
	msgProducer   *service.MessageProducer
}

func NewTxpoolApi (txpoolRepo *repository.TransactionRepository, messageProducer *service.MessageProducer) *TxpoolApi{
	txpConfig := conf.GetConfiguration().Txpool

	return &TxpoolApi{
		txRepository:  txpoolRepo,
		timeoutTicker: timeout.NewTimeoutTicker(txpConfig.TimeoutMs),
		maxTxByte:     txpConfig.MaxTransactionByte,
		msgProducer:   messageProducer,
	}
}

func (txpoolApi TxpoolApi) SaveTransaction(tx transaction.Transaction) error {
	if tx.TxStatus != transaction.VALID {
		return errors.New("transaction is not valid")
	}

	txData := tx.TxData
	hashArgs := []string{txData.Jsonrpc, string(txData.Method), string(txData.Params.Function), txData.ICodeID, tx.PublishPeerId, tx.TimeStamp.String(), string(tx.TxId), string(tx.TxType)}
	for _, str := range txData.Params.Args {
		hashArgs = append(hashArgs, str)
	}
	resultHash := common.ComputeSHA256(hashArgs)

	if resultHash != tx.TxHash {
		return errors.New("hash value is incorrect")
	}

	return (*txpoolApi.txRepository).Save(tx)
}

func (txpoolApi TxpoolApi) RemoveTransaction(transactionId transaction.TransactionId) error {
	return (*txpoolApi.txRepository).Remove(transactionId)
}