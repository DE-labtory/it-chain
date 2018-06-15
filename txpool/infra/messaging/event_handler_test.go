package messaging_test

import (
	"testing"

	"github.com/it-chain/it-chain-Engine/gateway/infra"
	"github.com/it-chain/it-chain-Engine/txpool"
	"github.com/it-chain/it-chain-Engine/txpool/infra/messaging"
)

type MockTransactionRepository struct {
	SaveFunc     func(transaction txpool.Transaction) error
	RemoveFunc   func(id txpool.TransactionId) error
	FindByIdFunc func(id txpool.TransactionId) (*txpool.Transaction, error)
	FindAllFunc  func() ([]*txpool.Transaction, error)
}

func (m MockTransactionRepository) Save(transaction txpool.Transaction) error {
	return m.SaveFunc(transaction)
}

func (m MockTransactionRepository) Remove(id txpool.TransactionId) error {
	return m.RemoveFunc(id)
}

func (m MockTransactionRepository) FindById(id txpool.TransactionId) (*txpool.Transaction, error) {
	return m.FindByIdFunc(id)
}

func (m MockTransactionRepository) FindAll() ([]*txpool.Transaction, error) {
	return m.FindAllFunc()
}

type MockLeaderRepository struct {
	GetLeaderFunc func() txpool.Leader
	SetLeaderFunc func(leader txpool.Leader)
}

func (m MockLeaderRepository) GetLeader() txpool.Leader {
	return m.GetLeaderFunc()
}

func (m MockLeaderRepository) SetLeader(leader txpool.Leader) {
	m.SetLeaderFunc(leader)
}

func TestTxEventHandler_HandleTxCreatedEvent(t *testing.T) {

	//given
	tests := map[string]struct {
		input  txpool.TxCreatedEvent
		output string
		err    error
	}{
		"handle success": {
			input: txpool.TxCreatedEvent{},
			err:   nil,
		},
		"handle fail": {
			input: txpool.TxCreatedEvent{},
			err:   infra.ErrConnAlreadyExist,
		},
	}

	mockTxRepo := MockTransactionRepository{}
	mockLeaderRepo := MockLeaderRepository{}

	event_handler := messaging.NewTxEventHandler(mockTxRepo, mockLeaderRepo)

	for testName, test := range tests {
		t.Logf("Running test case %s", testName)

		event_handler.HandleTxCreatedEvent()
	}

}
