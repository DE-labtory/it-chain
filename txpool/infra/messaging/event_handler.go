package messaging

import (
	"log"

	"github.com/it-chain/it-chain-Engine/txpool"
)

type TxEventHandler struct {
	txRepository     txpool.TransactionRepository
	leaderRepository txpool.LeaderRepository
}

//add tx to txrepository
func (t TxEventHandler) HandleTxCreatedEvent(txCreatedEvent txpool.TxCreatedEvent) {

	txID := txCreatedEvent.ID

	if txID == "" {
		return
	}

	tx := txCreatedEvent.GetTransaction()
	err := t.txRepository.Save(tx)

	if err != nil {
		log.Println(err.Error())
	}
}

//remove transaction
func (t TxEventHandler) HandleTxDeletedEvent(txDeletedEvent txpool.TxDeletedEvent) {

	txID := txDeletedEvent.ID

	if txID == "" {
		return
	}

	err := t.txRepository.Remove(txpool.TransactionId(txID))

	if err != nil {
		log.Println(err.Error())
	}
}

//update leadaer
func (t TxEventHandler) HandleLeaderChangedEvent(leaderChangedEvent txpool.LeaderChangedEvent) {

	leaderID := leaderChangedEvent.ID

	leader := txpool.Leader{
		txpool.LeaderId{leaderID},
	}

	t.leaderRepository.SetLeader(leader)
}
