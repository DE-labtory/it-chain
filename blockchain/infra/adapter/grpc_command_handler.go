package adapter

import "github.com/it-chain/it-chain-Engine/blockchain"

type GrpcCommandHandlerBlockApi interface {
	CreateGenesisBlock(genesisConfFilePath string) (blockchain.Block, error)
	CreateBlock(txList []blockchain.Transaction) (blockchain.Block, error)
	SyncedCheck(block blockchain.Block) error
}

type GrpcCommandHandler struct {
	blockApi GrpcCommandHandlerBlockApi
}

func NewGrpcCommandHandler(blockApi GrpcCommandHandlerBlockApi) *GrpcCommandHandler {
	return &GrpcCommandHandler{
		blockApi: blockApi,
	}
}

func (g *GrpcCommandHandler) HandleGrpcCommand(command blockchain.GrpcReceiveCommand) error {
	switch command.Protocol {
	case "SyncCheckRequestProtocol":
		//TODO: 상대방의 SyncCheck를 위해서 자신의 last block을 보내준다.
		break

	case "SyncCheckResponseProtocol":
		//TODO: 상대방의 last block을 받아서 SyncCheck를 시작한다.
		break

	case "BlockRequestProtocol":
		//TODO: Construct 과정을 위해서 상대방에게 block을 보내준다.
		break

	case "BlockResponseProtocol":
		//TODO: Construct 과정에서 block을 받는다.
		break
	}

	return nil
}
