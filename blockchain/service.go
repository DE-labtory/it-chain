package blockchain

import (
	"github.com/it-chain/it-chain-Engine/core/eventstore"
	"github.com/it-chain/midgard"
)

type BlockQueryService interface {
	BlockQueryInnerService
	BlockQueryOuterService
}

type BlockQueryInnerService interface {
	GetLastBlock() (Block, error)
}

type BlockQueryOuterService interface {
	GetLastBlockFromPeer(peer Peer) (Block, error)
	GetBlockByHeightFromPeer(peer Peer, height BlockHeight) (Block, error)
}

type PeerService interface {
	GetRandomPeer() (Peer, error)
}

func CommitBlock(block Block) error {

	event, err := createBlockCommittedEvent(block)

	if err != nil {
		return err
	}

	blockId := string(block.GetSeal())
	eventstore.Save(blockId, event)

	return nil
}

func createBlockCommittedEvent(block Block) (*BlockCommittedEvent, error) {

	aggregateId := string(block.GetSeal())

	return &BlockCommittedEvent{
		EventModel: midgard.EventModel{
			ID: aggregateId,
		},
		State: Committed,
	}, nil
}
