package service

import (
	"it-chain/db/leveldbhelper"
	"it-chain/domain"
	"it-chain/common"
	"it-chain/network/comm"
	pb "it-chain/network/protos"
)

const (
	WAITING_TRANSACTION = "waiting_transaction"
)

type TransactionServiceImpl struct {
	DB *leveldbhelper.DBProvider
	Comm comm.ConnectionManager
	PeerService PeerService
}

func CreateNewTransactionService(path string, comm comm.ConnectionManager, ps PeerService) *TransactionServiceImpl {
	return &TransactionServiceImpl{DB: leveldbhelper.CreateNewDBProvider(path), Comm: comm, PeerService: ps}
}

func (t *TransactionServiceImpl) Close() {
	t.DB.Close()
}

func (t *TransactionServiceImpl) AddTransaction(tx *domain.Transaction) error {
	db := t.DB.GetDBHandle(WAITING_TRANSACTION)
	serializedTX, err := common.Serialize(tx)
	if err != nil {
		return err
	}

	err = db.Put([]byte(tx.TransactionID), serializedTX, true)
	return err
}

func (t *TransactionServiceImpl) DeleteTransactions(txs []*domain.Transaction) error {
	db := t.DB.GetDBHandle(WAITING_TRANSACTION)
	batch := make(map[string][]byte)

	for _, tx := range txs {
		batch[tx.TransactionID] = nil
	}

	return db.WriteBatch(batch, true)
}

func (t *TransactionServiceImpl) GetTransactions(limit int) ([]*domain.Transaction, error) {
	db := t.DB.GetDBHandle(WAITING_TRANSACTION)
	iter := db.GetIteratorWithPrefix()
	ret := make([]*domain.Transaction, 0)
	cnt := 0

	for iter.Next() {
		val := iter.Value()
		tx := &domain.Transaction{}
		err := common.Deserialize(val, tx)

		if err != nil {
			return nil, err
		}

		ret = append(ret, tx)
		cnt++
		if cnt == limit {
			break
		}
	}

	iter.Release()

	return ret, nil
}

func (t *TransactionServiceImpl) SendToLeader(interface{}) {
	txs, err := t.GetTransactions(1)
	if err != nil {
		logger.Println("Error on GetTransactions")
	}

	message := &pb.StreamMessage{}
	message.Content = &pb.StreamMessage_Transaction{
		Transaction: &pb.Transaction{},
	}

	if err !=nil{
		logger.Println("fail to serialize message")
	}

	errorCallBack := func(onError error) {
		logger.Println("fail to send message error:", onError.Error())
	}

	t.Comm.SendStream(message, errorCallBack, t.PeerService.GetLeader().PeerID)
	t.DeleteTransactions(txs)
}