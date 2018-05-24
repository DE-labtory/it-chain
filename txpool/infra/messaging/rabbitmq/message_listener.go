package rabbitmq

import (
	"encoding/json"

	"github.com/it-chain/it-chain-Engine/messaging/rabbitmq/event"
	"github.com/it-chain/it-chain-Engine/txpool/api"
	"github.com/streadway/amqp"
	"github.com/it-chain/it-chain-Engine/txpool/domain/model/transaction"
)

type MessageListener struct {
	txpoolApi api.TxpoolApi
}

func (ml MessageListener) ListenTransactionCreateEvent(messageChannel <-chan amqp.Delivery) {
	go func() {
		for message := range messageChannel {
			eventMessage := &event.TransactionCreateEvent{}
			err := json.Unmarshal(message.Body, &eventMessage)

			if err != nil {
				//error
				// TODO 에러처리하기
			}

			var txDataMessage transaction.TxData
			json.Unmarshal(eventMessage.TransactionData, &txDataMessage)
			ml.txpoolApi.SaveTxData(eventMessage.PeerId, transaction.TxDataType(eventMessage.TxDataType), txDataMessage)
		}
	}()
}

func (ml MessageListener) ListenTransactionReceiveEvent(messageChannel <-chan amqp.Delivery) {
	go func() {
		for message := range messageChannel {
			eventMessage := &event.TransactionReceiveEvent{}
			err := json.Unmarshal(message.Body, &eventMessage)

			if err != nil {
				//error
				// TODO 에러처리하기
			}

			var txMessage transaction.Transaction
			json.Unmarshal(eventMessage.Transaction, &txMessage)
			ml.txpoolApi.SaveTransaction(txMessage)
		}
	}()
}

func (ml MessageListener) ListenLeaderChangeEvent(messageChannel <-chan amqp.Delivery) {
	go func() {
		for message := range messageChannel {
			eventMessage := &event.LeaderChangeEvent{}
			err := json.Unmarshal(message.Body, &eventMessage)

			if err != nil {
				//error
				// TODO 에러처리하기
			}

			//TODO event 메세지 처리하기
		}
	}()
}
