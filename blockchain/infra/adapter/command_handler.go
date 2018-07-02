package adapter

import (
	"github.com/it-chain/it-chain-Engine/blockchain"
	"github.com/it-chain/it-chain-Engine/txpool"
	"errors"
	"github.com/it-chain/midgard"
)

var ErrBlockNil = errors.New("Block nil error");

type BlockApi interface {
	AddBlockToPool(block blockchain.Block)
}

type CommandHandler struct {
	blockApi BlockApi
	repositoryProjector RepositoryProjector
}

func NewCommandHandler(blockApi BlockApi, repositoryProjector RepositoryProjector) *CommandHandler {
	return &CommandHandler{
		blockApi: blockApi,
		repositoryProjector: repositoryProjector,
	}
}

// txpool에서 받은 transactions들을 block으로 만들어서 consensus에 보내준다.
func (h *CommandHandler) HandleProposeBlockCommand(command blockchain.ProposeBlockCommand) {
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

/// 합의된 block이 넘어오면 block pool에 저장한다.
func (h *CommandHandler) HandleConfirmBlockCommand(command blockchain.ConfirmBlockCommand) error {
	block := command.Block
	if block == nil {
		return ErrBlockNil
	}

	h.blockApi.AddBlockToPool(block)

	// BlockQueuedEvent eventRepository에 추가하기
	event := createBlockQueuedEvent(block)
	h.repositoryProjector.EventRepository.Save(blockchain.BLOCK_QUEUED_EID, event)
	return nil
}

func createBlockQueuedEvent(block blockchain.Block) blockchain.BlockQueuedEvent {
	return blockchain.BlockQueuedEvent{
		EventModel: midgard.EventModel{
			ID: blockchain.BLOCK_QUEUED_EID,
		},
		Block: block,
	}
}
