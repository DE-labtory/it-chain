package adapter

import (
	"log"

	"github.com/it-chain/engine/txpool"
	"github.com/it-chain/engine/txpool/api"
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

	txData := txpool.TxData{
		ICodeID: txCreateCommand.ICodeID,
		Jsonrpc: txCreateCommand.Jsonrpc,
		Method:  txCreateCommand.Method,
		Params:  txCreateCommand.Params,
	}

	tx, err := t.transactionApi.CreateTransaction(txData)

	if err != nil {
		log.Println(err.Error())
	}

	log.Printf("transactions are created [%v]", tx)
}
