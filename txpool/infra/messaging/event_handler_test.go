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

	tests := map[string]struct {
		input txpool.TxCreatedEvent
		output interface{}
		err error
	} {
		"success": {
			input: txpool.TxCreatedEvent{
				EventModel: midgard.EventModel{
					ID: "zf",
				},
				TimeStamp: time.Now().UTC(),
				TxHash: "zf_hash",
				TxData: txpool.TxData{
					ID: "zf2",
				},
			},
			output: nil,
			err: nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		// When
		teh.HandleTxCreatedEvent(test.input)

		// Then
		tx, _ := teh.txRepository.FindById("zf")
		assert.Equal(t, (*tx).TxId, txpool.TransactionId("zf"))
	}
}

func TestTxEventHandler_HandleTxDeletedEvent(t *testing.T) {
	dbPath := "./test"

	defer os.RemoveAll(dbPath)

	teh := TxEventHandler{
		txRepository: leveldb.NewTransactionRepository(dbPath),
		leaderRepository: txpool.NewLeaderRepository(),
	}

	tests := map[string]struct {
		input txpool.TxDeletedEvent
		output interface{}
		err error
	} {
		"success": {
			input: txpool.TxDeletedEvent{
				EventModel: midgard.EventModel{
					ID: "zf",
				},
			},
			output: nil,
			err: nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		// When
		teh.txRepository.Save(txpool.Transaction{
			TxId: "zf",
			TimeStamp: time.Now().UTC(),
			TxData: txpool.TxData{
				ID: "zf2",
			},
		})
		teh.HandleTxDeletedEvent(test.input)

		// Then
		txs, err := teh.txRepository.FindAll()
		assert.Equal(t, 0, len(txs))
		assert.Equal(t, nil, err)
	}
}

func TestTxEventHandler_HandleLeaderChangedEvent(t *testing.T) {
	dbPath := "./test"

	defer os.RemoveAll(dbPath)

	teh := TxEventHandler{
		txRepository: leveldb.NewTransactionRepository(dbPath),
		leaderRepository: txpool.NewLeaderRepository(),
	}

	tests := map[string]struct {
		input txpool.LeaderChangedEvent
		output interface{}
		err error
	} {
		"success": {
			input: txpool.LeaderChangedEvent{
				EventModel: midgard.EventModel{
					ID: "zf",
				},
			},
			output: nil,
			err: nil,
		},
	}

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		// When
		teh.HandleLeaderChangedEvent(test.input)

		// Then
		leader := teh.leaderRepository.GetLeader()
		assert.Equal(t, "zf", leader.GetID())
	}
}