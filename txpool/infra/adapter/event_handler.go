package adapter

import (
	"errors"

	"github.com/it-chain/it-chain-Engine/txpool"
)

var ErrNoEventID = errors.New("no event_store id ")

//////////////Event Handler

type EventHandler struct {
}

//RepositoryProjector
//do not import any api or service
//event를 받아서 repository를 update하는 역할만 수행
type RepositoryProjector struct {
	txRepository     txpool.TransactionRepository
	leaderRepository txpool.LeaderRepository
}

func NewTxEventHandler(txRepository txpool.TransactionRepository, leaderRepository txpool.LeaderRepository) *RepositoryProjector {
	return &RepositoryProjector{
		txRepository:     txRepository,
		leaderRepository: leaderRepository,
	}
}

//add tx to txrepository
func (t RepositoryProjector) HandleTxCreatedEvent(txCreatedEvent txpool.TxCreatedEvent) error {

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
func (t RepositoryProjector) HandleTxDeletedEvent(txDeletedEvent txpool.TxDeletedEvent) error {

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
func (t RepositoryProjector) HandleLeaderChangedEvent(leaderChangedEvent txpool.LeaderChangedEvent) error {

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
