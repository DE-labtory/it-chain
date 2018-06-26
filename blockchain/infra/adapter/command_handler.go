package adapter

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/txpool"
)

type CommandHandlerNodeApi interface {
	UpdateNode(node blockchain.Node) error
}

type CommandHandlerBlockApi interface {
	CreateGenesisBlock(genesisConfFilePath string) (blockchain.Block, error)
	CreateBlock(txList []blockchain.Transaction) (blockchain.Block, error)
}

type CommandHandler struct {
	nodeApi CommandHandlerNodeApi
	blockApi CommandHandlerBlockApi
}

func NewCommandHandler(blockApi CommandHandlerBlockApi, nodeApi CommandHandlerNodeApi) *CommandHandler {
	return &CommandHandler{
		nodeApi: nodeApi,
		blockApi: blockApi,
	}
}

/// 임의로 선정한 노드의 정보를 업데이트한다.
func (b *CommandHandler) HandleUpdateNodeCommand(command blockchain.NodeUpdateCommand) {
	//eventID := command.GetID()
	//
	//if eventID == "" {
	//	log.Println(ErrEmptyEventId)
	//	return
	//}
	//
	//node := command.Node
	//
	//if node.IpAddress == "" {
	//	log.Println(ErrEmptyIpAddress)
	//	return
	//}
	//
	//if node.NodeId.Id == "" {
	//	log.Println(ErrEmptyNodeId)
	//	return
	//}
	//
	//err := b.nodeApi.UpdateNode(node)
	//
	//if err != nil {
	//	log.Println("%s: %s", ErrNodeApi, err)
	//}
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
