package blockchain

import "github.com/it-chain/it-chain-Engine/core/eventstore"

func CommitBlock(block Block) error {
	event, err := createBlockCommittedEvent(block)
	if err != nil {
		return err
	}
	blockId := string(block.GetSeal())
	eventstore.Save(blockId, event)
	return nil
}
