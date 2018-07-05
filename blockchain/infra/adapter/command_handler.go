package adapter

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/blockchain/api"
	"github.com/it-chain/it-chain-Engine/txpool"
)

type BlockchainCommandHandler struct {
	nodeApi  api.NodeApi
	blockApi api.BlockApi
}

func NewBlockchainCommandHandler(blockApi api.BlockApi, nodeApi api.NodeApi) *BlockchainCommandHandler {
	return &BlockchainCommandHandler{
		nodeApi:  nodeApi,
		blockApi: blockApi,
	}
}

type CommandHandlerBlockApi interface {
	CreateGenesisBlock(genesisConfFilePath string) (blockchain.Block, error)
	CreateBlock(txList []blockchain.Transaction) (blockchain.Block, error)
}

type CommandHandler struct {
	blockApi CommandHandlerBlockApi
}

func NewCommandHandler(blockApi CommandHandlerBlockApi) *CommandHandler {
	return &CommandHandler{
		blockApi: blockApi,
	}
}

func (handler *CommandHandler) HandleProposeBlockCommand(command blockchain.ProposeBlockCommand) {
	//rawTxList := command.Transactions
	//
	//txList, err := convertTxList(rawTxList)
	//if err != nil {
	//	// TODO: handle errors
	//	return
	//}
	//
	//block, err := handler.blockApi.CreateBlock(txList)
	//if err != nil {
	//	// TODO: handle errors
	//	return
	//}
	// TODO: service는 api에서 호출되어야한다.
	//dispatcher.SendBlockValidateCommand(block)
}

// TODO: yggdrasill/impl/Transaction과 txpool/Transaction이 다름.
func convertTxList(txList []txpool.Transaction) ([]blockchain.Transaction, error) {
	return nil, nil
}
