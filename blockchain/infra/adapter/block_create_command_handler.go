package adapter

import (
	"log"

	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/txpool"
)

type CreateBlockApi interface {
	CreateBlock(txList []blockchain.Transaction) error
}

type BlockCreateCommandHandler struct {
	blockApi CreateBlockApi
}

func NewBlockCreateCommandHandler(blockApi CreateBlockApi) *BlockCreateCommandHandler {
	return &BlockCreateCommandHandler{
		blockApi: blockApi,
	}
}

// txpool에서 받은 transactions들을 block으로 만들어서 consensus에 보내준다.
// TODO: write test code
func (h *BlockCreateCommandHandler) HandleProposeBlockCommand(command blockchain.ProposeBlockCommand) {
	txList := ConvertTxList(command.Transactions)

	err := h.blockApi.CreateBlock(txList)
	if err != nil {
		log.Fatal(err)
	}
}

func ConvertTxList(txpoolTxList []txpool.Transaction) []blockchain.Transaction {
	txList := make([]blockchain.Transaction, 0)
	for _, tx := range txpoolTxList {
		txList = append(txList, convertTx(tx))
	}
	return txList
}

func convertTx(tx txpool.Transaction) *blockchain.DefaultTransaction {
	return &blockchain.DefaultTransaction{
		ID:        tx.TxId,
		Status:    blockchain.StatusTransactionInvalid,
		PeerID:    tx.PublishPeerId,
		Timestamp: tx.TimeStamp,
		TxData: &blockchain.TxData{
			Jsonrpc: tx.TxData.Jsonrpc,
			Method:  blockchain.TxDataType(tx.TxData.Method),
			Params: blockchain.Params{
				Function: tx.TxData.Params.Function,
				Args:     tx.TxData.Params.Args,
			},
			ID: tx.TxData.ID,
		},
		Signature: tx.Signature,
	}
}
