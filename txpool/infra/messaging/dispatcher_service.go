package messaging

import (
	"github.com/it-chain/it-chain-Engine/txpool"
)

type Publisher func(exchange string, topic string, data interface{}) (err error) 	//해당 publish함수는 midgard 에서 의존성 주입을 받기 위해 interface로 작성한다.
																						//모든 의존성 주입은 컴포넌트.go 에서 이루어짐

//todo implement create command using transaction and leader and send to rabbitmq
type MessageDispatcher struct {
	publisher Publisher // midgard.client
}

//todo implement sendTransactionCommand 정의 해야함
func (m MessageDispatcher) SendTransactions(transactions []txpool.Transaction, leader txpool.Leader) error {
	return nil
}

//todo implement proposeBlockCommand 정의 해야함
func (m MessageDispatcher) ProposeBlock(transactions []txpool.Transaction) error {
	return nil
}
