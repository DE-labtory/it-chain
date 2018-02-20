package service

import (
	"it-chain/db/leveldbhelper"
	"it-chain/domain"
	"it-chain/common"
	"it-chain/network/comm"
	pb "it-chain/network/protos"
	"github.com/spf13/viper"
	"strconv"
	"time"
	"github.com/rs/xid"
	"github.com/pkg/errors"
)

const (
	WAITING_TRANSACTION = "waiting_transaction"
)

type TransactionServiceImpl struct {
	DB *leveldbhelper.DBProvider
	Comm comm.ConnectionManager
	PeerService PeerService
}

func NewTransactionService(path string, comm comm.ConnectionManager, ps PeerService) *TransactionServiceImpl {
	transactionService := &TransactionServiceImpl{DB: leveldbhelper.CreateNewDBProvider(path), Comm: comm, PeerService: ps}

	i, _ := strconv.Atoi(viper.GetString("batchTimer.pushPeerTable"))

	broadCastPeerTableBatcher := NewBatchService(time.Duration(i)*time.Second,transactionService.SendToLeader,false)
	broadCastPeerTableBatcher.Add("Send tx to leader")
	broadCastPeerTableBatcher.Start()

	//comm.Subscribe()

	return transactionService
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
	
	//todo max 몇개까지 보낼것인지
	txs, err := t.GetTransactions(100)

	if err != nil {
		common.Log.Println("Error on GetTransactions")
	}

	if len(txs) == 0{
		common.Log.Println("No transactions to send")
		return
	}

	message := &pb.StreamMessage{}
	message.Content = &pb.StreamMessage_Transaction{
		Transaction: &pb.Transaction{},
	}

	if err !=nil{
		common.Log.Println("fail to serialize message")
	}

	errorCallBack := func(onError error) {
		common.Log.Println("fail to send message error:", onError.Error())
	}

	if t.PeerService.GetLeader() != nil {
		t.Comm.SendStream(message, errorCallBack, t.PeerService.GetLeader().PeerID)
		t.DeleteTransactions(txs)
	}
}

func (t *TransactionServiceImpl) CreateTransaction(txData *domain.TxData) (*domain.Transaction, error){

	transaction := domain.CreateNewTransaction(
		t.PeerService.GetPeerTable().MyID,
		xid.New().String(),
		domain.General,
		time.Now(),
		txData)

	err := t.AddTransaction(transaction)

	if err != nil{
		return nil, errors.New("faild to add transaction")
	}

	return transaction, nil
}