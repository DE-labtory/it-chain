package service

import (
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/messaging/event"
	"github.com/it-chain/it-chain-Engine/txpool/domain/model"
	"github.com/it-chain/it-chain-Engine/txpool/domain/model/transaction"
)

type Publish func(topic string, data []byte) error

type MessageService struct {
	Publish Publish
}

func NewMessageApi(publish Publish) *MessageService {
	return &MessageService{Publish: publish}
}

func (mApi *MessageService) SendTransaction(transaction transaction.Transaction, leader model.Leader) error {
	txData, err := transaction.Serialize()
	if err != nil {
		return err
	}
	common.Serialize(event.TransactionSendEvent{
		LeaderId:    leader.GetStringID(),
		Transaction: txData,
	})
	return nil
}
