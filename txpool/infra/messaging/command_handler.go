package messaging

import (
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/it-chain-Engine/txpool/api"
)

type TxCommandHandler struct {
	txpoolApi api.TxpoolApi
}

func (t TxCommandHandler) HandleTxCreate(txCreateCommand txpool.TxCreateCommand) {

}
