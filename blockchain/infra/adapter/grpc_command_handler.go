package adapter

import (
	"errors"

	"encoding/json"

	"github.com/it-chain/it-chain-Engine/blockchain"
)

var ErrBlockInfoDeliver = errors.New("block info deliver failed")
var ErrGetBlock = errors.New("error when Getting block")
var ErrResponseBlock = errors.New("error when response block")
var ErrGetLastBlock = errors.New("error when get last block")
var ErrSyncCheckResponse = errors.New("error when sync check response")

type BlockApi interface {
	SyncedCheck(block blockchain.Block) error
}

type ReadOnlyBlockRepository interface {
	NewEmptyBlock() (blockchain.Block, error)
	GetLastBlock(block blockchain.Block) error
	GetBlockByHeight(block blockchain.Block, height uint64) error
}

type SyncGrpcCommandService interface {
	ResponseBlock(peerId blockchain.PeerId, block blockchain.Block) error
	SyncCheckResponse(block blockchain.Block) error
}

type GrpcCommandHandler struct {
	blockApi           BlockApi
	blockRepository    ReadOnlyBlockRepository
	grpcCommandService SyncGrpcCommandService
}

func NewGrpcCommandHandler(blockApi BlockApi, blockRepository ReadOnlyBlockRepository, grpcCommandService SyncGrpcCommandService) *GrpcCommandHandler {
	return &GrpcCommandHandler{
		blockApi:           blockApi,
		blockRepository:    blockRepository,
		grpcCommandService: grpcCommandService,
	}
}

func (g *GrpcCommandHandler) HandleGrpcCommand(command blockchain.GrpcReceiveCommand) error {
	switch command.Protocol {
	case "SyncCheckRequestProtocol":
		//TODO: 상대방의 SyncCheck를 위해서 자신의 last block을 보내준다.
		block := &blockchain.DefaultBlock{}

		err := g.blockRepository.GetLastBlock(block)
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

		block := &blockchain.DefaultBlock{}

		var height uint64
		err := json.Unmarshal(command.Body, &height)
		if err != nil {
			return ErrBlockInfoDeliver
		}

		err = g.blockRepository.GetBlockByHeight(block, height)
		if err != nil {
			return ErrGetBlock
		}

		err = g.grpcCommandService.ResponseBlock(command.FromPeer.PeerId, block)
		if err != nil {
			return ErrResponseBlock
		}
		break

	case "BlockResponseProtocol":
		//TODO: Construct 과정에서 block을 받는다.
		break
	}

	return nil
}
