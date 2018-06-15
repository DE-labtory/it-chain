package messaging

import (
	"testing"
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/midgard"
	"time"
	"github.com/it-chain/it-chain-Engine/txpool/infra/repository/leveldb"
	"os"
	"github.com/stretchr/testify/assert"
)

func TestTxEventHandler_HandleTxCreatedEvent(t *testing.T) {
	dbPath := "./test"

	defer os.RemoveAll(dbPath)

	teh := TxEventHandler{
		txRepository: leveldb.NewTransactionRepository(dbPath),
		leaderRepository: txpool.NewLeaderRepository(),
	}

	// When
	event := txpool.TxCreatedEvent{
		EventModel: midgard.EventModel{
			ID: "zf",
		},
		TimeStamp: time.Now().UTC(),
		TxHash:    "zf_hash",
		TxData: txpool.TxData{
			ID: "zf2",
		},
	}

	teh.HandleTxCreatedEvent(event)

	// Then
	tx, _ := teh.txRepository.FindById("zf")
	assert.Equal(t, (*tx).TxId, txpool.TransactionId("zf"))
}

func TestTxEventHandler_HandleTxCreatedEvent_WhenTxIdMissing(t *testing.T) {
	dbPath := "./test"

	defer os.RemoveAll(dbPath)

	teh := TxEventHandler{
		txRepository: leveldb.NewTransactionRepository(dbPath),
		leaderRepository: txpool.NewLeaderRepository(),
	}

	// When
	event := txpool.TxCreatedEvent{
		EventModel: midgard.EventModel{
			// ID Missing
		},
		TimeStamp: time.Now().UTC(),
		TxHash:    "zf_hash",
		TxData: txpool.TxData{
			ID: "zf2",
		},
	}

	teh.HandleTxCreatedEvent(event)

	// Then
	tx, _ := teh.txRepository.FindById("zf")
	assert.Equal(t, (*txpool.Transaction)(nil) , tx)
}

func TestTxEventHandler_HandleTxDeletedEvent(t *testing.T) {
	dbPath := "./test"

	defer os.RemoveAll(dbPath)

	// When
	teh := TxEventHandler{
		txRepository: leveldb.NewTransactionRepository(dbPath),
		leaderRepository: txpool.NewLeaderRepository(),
	}

	event := txpool.TxDeletedEvent{
		EventModel: midgard.EventModel{
			ID: "zf",
		},
	}

	teh.HandleTxDeletedEvent(event)

	// Then
	tx, err := teh.txRepository.FindById("zf")
	assert.Equal(t, (*txpool.Transaction)(nil) , tx)
	assert.Equal(t, nil, err)
}

func TestTxEventHandler_HandleTxDeletedEvent_WhenTxIdMissing(t *testing.T) {
	dbPath := "./test"

	defer os.RemoveAll(dbPath)

	// When
	teh := TxEventHandler{
		txRepository: leveldb.NewTransactionRepository(dbPath),
		leaderRepository: txpool.NewLeaderRepository(),
	}

	event := txpool.TxDeletedEvent{
		EventModel: midgard.EventModel{
			// ID Missing
		},
	}

	teh.HandleTxDeletedEvent(event)

	// Then
	tx, err := teh.txRepository.FindById("zf")
	assert.Equal(t, (*txpool.Transaction)(nil) , tx)
	assert.Equal(t, nil, err)
}

func TestTxEventHandler_HandleLeaderChangedEvent(t *testing.T) {
	dbPath := "./test"

	defer os.RemoveAll(dbPath)

	// When
	teh := TxEventHandler{
		txRepository: leveldb.NewTransactionRepository(dbPath),
		leaderRepository: txpool.NewLeaderRepository(),
	}

	event := txpool.LeaderChangedEvent{
		EventModel: midgard.EventModel{
			ID: "zf",
		},
	}

	teh.HandleLeaderChangedEvent(event)

	// Then
	leader := teh.leaderRepository.GetLeader()
	assert.Equal(t, "zf", leader.GetID())
}