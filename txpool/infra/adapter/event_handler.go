package adapter

import (
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/pkg/errors"
)

var ErrNoEventID = errors.New("no event id ")

type TxEventHandler struct {
	txRepository     txpool.TransactionRepository
	leaderRepository txpool.LeaderRepository
}

func NewTxEventHandler(txRepository txpool.TransactionRepository, leaderRepository txpool.LeaderRepository) *TxEventHandler {
	return &TxEventHandler{
		txRepository:     txRepository,
		leaderRepository: leaderRepository,
	}
}

//add tx to txrepository
func (t TxEventHandler) HandleTxCreatedEvent(txCreatedEvent txpool.TxCreatedEvent) error {

	txID := txCreatedEvent.ID

	if txID == "" {
		return ErrNoEventID
	}

	tx := txCreatedEvent.GetTransaction()
	err := t.txRepository.Save(tx)

	if err != nil {
		return err
	}

	return nil
}

//remove transaction
func (t TxEventHandler) HandleTxDeletedEvent(txDeletedEvent txpool.TxDeletedEvent) error {

	txID := txDeletedEvent.ID

	if txID == "" {
		return ErrNoEventID
	}

	err := t.txRepository.Remove(txpool.TransactionId(txID))

	if err != nil {
		return err
	}

	return nil
}

//update leader
func (t TxEventHandler) HandleLeaderChangedEvent(leaderChangedEvent txpool.LeaderChangedEvent) error {

	leaderID := leaderChangedEvent.ID

	if leaderID == "" {
		return ErrNoEventID
	}

	leader := txpool.Leader{
		txpool.LeaderId{leaderID},
	}

	t.leaderRepository.SetLeader(leader)

	return nil
}
