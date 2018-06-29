package adapter

import (
	"log"

	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/it-chain-Engine/txpool/api"
)

type TxCommandHandler struct {
	transactionApi api.TransactionApi
}

func NewTxCommandHandler(transactionApi api.TransactionApi) *TxCommandHandler {
	return &TxCommandHandler{
		transactionApi: transactionApi,
	}
}

func (t *TxCommandHandler) HandleTxCreateCommand(txCreateCommand txpool.TxCreateCommand) {

	txID := txCreateCommand.GetID()

	if txID == "" {
		log.Println("need id")
		return
	}

	txData := txCreateCommand.TxData
	err := t.transactionApi.CreateTransaction(txID, txData)

	if err != nil {
		log.Println(err.Error())
	}

}
