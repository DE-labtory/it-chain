package api

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/conf"
	"github.com/it-chain/it-chain-Engine/txpool"
)

type TxpoolApi struct {
	txRepository  *txpool.TransactionRepository
	timeoutTicker *txpool.TimeoutTicker
	maxTxByte     int
	msgProducer   *txpool.MessageProducer
}

func NewTxpoolApi(txpoolRepo *txpool.TransactionRepository, messageProducer *txpool.MessageProducer) *TxpoolApi {
	txpConfig := conf.GetConfiguration().Txpool

	return &TxpoolApi{
		txRepository:  txpoolRepo,
		timeoutTicker: txpool.NewTimeoutTicker(txpConfig.TimeoutMs),
		maxTxByte:     txpConfig.MaxTransactionByte,
		msgProducer:   messageProducer,
	}
}

func (txpoolApi TxpoolApi) SaveTransaction(tx txpool.Transaction) error {
	if tx.TxStatus != txpool.VALID {
		return errors.New("transaction is not valid")
	}

	if tx.TxHash != tx.CalcHash() {
		return errors.New("hash value is incorrect")
	}

	return (*txpoolApi.txRepository).Save(tx)
}

func (txpoolApi TxpoolApi) SaveTxData(publishPeerId string, txType txpool.TxDataType, txData txpool.TxData) error {
	tx := txpool.NewTransaction(publishPeerId, txType, &txData)

	return txpoolApi.SaveTransaction(*tx)
}

func (txpoolApi TxpoolApi) RemoveTransaction(transactionId txpool.TransactionId) error {
	return (*txpoolApi.txRepository).Remove(transactionId)
}
