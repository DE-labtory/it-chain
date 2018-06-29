package adapter

import (
	"errors"

	"encoding/json"

	"github.com/it-chain/it-chain-Engine/blockchain"
)

var ErrBlockInfoDeliver = errors.New("block info deliver failed")
var ErrGetBlock = errors.New("error when Getting block")
var ErrResponseBlock = errors.New("error when response block")

type BlockApi interface {
	SyncedCheck(block blockchain.Block) error
}

type ReadOnlyBlockRepository interface {
	NewEmptyBlock() (blockchain.Block, error)
	GetLastBlock(block blockchain.Block) error
	GetBlockByHeight(blockHeight uint64) (blockchain.Block, error)
}

type GrpcCommandHandler struct {
	blockApi           BlockApi
	blockRepository    ReadOnlyBlockRepository
	grpcCommandService GrpcCommandService
}

func NewGrpcCommandHandler(blockApi BlockApi) *GrpcCommandHandler {
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
		var height uint64
		err := json.Unmarshal(command.Body, &height)
		if err != nil {
			return ErrBlockInfoDeliver
		}

		block, err := g.blockRepository.GetBlockByHeight(height)
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
