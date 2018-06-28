package adapter

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/p2p"
)

type BlockApi interface {
	SyncedCheck(block blockchain.Block) error
}

type ReadOnlyBlockRepository interface {
	GetLastBlock(block blockchain.Block) error
}

type SyncCheckGrpcCommandService interface {
	SyncCheckResponse(peerId p2p.PeerId, )
}

type GrpcCommandHandler struct {
	blockApi BlockApi
	blockRepository ReadOnlyBlockRepository
}

func NewGrpcCommandHandler(blockApi BlockApi, blockRepository ReadOnlyBlockRepository) *GrpcCommandHandler {
	return &GrpcCommandHandler{
		blockApi: blockApi,
		blockRepository: blockRepository,
	}
}

func (g *GrpcCommandHandler) HandleGrpcCommand(command blockchain.GrpcReceiveCommand) error {
	switch command.Protocol {
	case "SyncCheckRequestProtocol":
		//TODO: 상대방의 SyncCheck를 위해서 자신의 last block을 보내준다.
		block := blockchain.DefaultBlock{}
		g.blockRepository.GetLastBlock(&block)


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