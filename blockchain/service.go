package blockchain

import "github.com/it-chain/it-chain-Engine/core/eventstore"

type BlockService interface {
	GetLastBlock() (Block, error)
}

type PeerService interface {
	PeerInnerService
	PeerOuterService
}

type PeerInnerService interface {
	GetRandomPeer() (Peer, error)
}

type PeerOuterService interface {
	GetLastBlock(peer Peer) (Block, error)
	GetBlockByHeight(peer Peer, height BlockHeight) (Block, error)
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
