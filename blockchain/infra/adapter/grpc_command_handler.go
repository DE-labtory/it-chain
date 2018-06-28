package adapter

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"errors"
)

var ErrCreateBlock = errors.New("error when creating block")
var ErrGetLastBlock = errors.New("error when get last block")
var ErrSyncCheckResponse = errors.New("error when sync check response")

type BlockApi interface {
	SyncedCheck(block blockchain.Block) error
}

type ReadOnlyBlockRepository interface {
	NewEmptyBlock() (blockchain.Block, error)
	GetLastBlock(block blockchain.Block) error
}

type SyncCheckGrpcCommandService interface {
	SyncCheckResponse(block blockchain.Block) error
}

type GrpcCommandHandler struct {
	blockApi BlockApi
	blockRepository ReadOnlyBlockRepository
	grpcCommandService SyncCheckGrpcCommandService
}

func NewGrpcCommandHandler(blockApi BlockApi, blockRepository ReadOnlyBlockRepository, grpcCommandService SyncCheckGrpcCommandService) *GrpcCommandHandler {
	return &GrpcCommandHandler{
		blockApi: blockApi,
		blockRepository: blockRepository,
		grpcCommandService: grpcCommandService,
	}
}

func (g *GrpcCommandHandler) HandleGrpcCommand(command blockchain.GrpcReceiveCommand) error {
	switch command.Protocol {
	case "SyncCheckRequestProtocol":
		//TODO: 상대방의 SyncCheck를 위해서 자신의 last block을 보내준다.
		block, err := g.blockRepository.NewEmptyBlock()
		if err != nil {
			return ErrCreateBlock
		}

		err = g.blockRepository.GetLastBlock(block)
		if err != nil {
			return ErrGetLastBlock
		}

		err = g.grpcCommandService.SyncCheckResponse(block)
		if err != nil {
			return ErrSyncCheckResponse
		}
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