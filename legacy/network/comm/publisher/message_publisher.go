package publisher

import (
	"github.com/it-chain/it-chain-Engine/common"
	"github.com/it-chain/it-chain-Engine/legacy/network/comm/msg"
	"errors"
	"github.com/it-chain/it-chain-Engine/legacy/auth"
	"sync"
)

var logger_event_publisher = common.GetLogger("message_publisher.go")

//message를 받으면 관심이 있는 subscriber에게 전달하는 역활을 한다.
type MessagePublisher struct {
	subscribers map[string]func(message msg.OutterMessage)
	crpyto      auth.Crypto
	lock        *sync.Mutex
}

func NewMessagePublisher(crpyto auth.Crypto) *MessagePublisher{

	return &MessagePublisher{
		subscribers: make(map[string]func(message msg.OutterMessage)),
		crpyto: crpyto,
		lock: &sync.Mutex{},
	}
}

//subscriber를 등록한다.
func (mp *MessagePublisher) AddSubscriber(name string, subfunc func(message msg.OutterMessage)) error{

	mp.lock.Lock()
	defer mp.lock.Unlock()

	_, ok := mp.subscribers[name]

	if ok{
		return errors.New("already subscribed function")
	}

	mp.subscribers[name] = subfunc

	return nil
}

func (mp *MessagePublisher) ReceivedMessageHandle(message msg.OutterMessage){

	defer recover()

	envelope := message.Envelope

	if envelope == nil{
		logger_event_publisher.Info("message is nil", message)
		return
	}

	//todo Verify 부분 추가해야함
	//_, pub, err := mp.crpyto.GetKey()

	//if err != nil{
	//	logger_event_publisher.Infof("failed to load key: %s", err.Error())
	//	return
	//}

	//vaild, err := mp.crpyto.Verify(pub,envelope.Signature,envelope.Payload,nil)
	//
	//if !vaild || err !=nil{
	//	logger_event_publisher.Info("failed to verify message")
	//	return
	//}

	msg, err := envelope.GetMessage()

	message.Message = msg

	if err != nil{
		logger_event_publisher.Info("failed to Unmarshal message:", message)
		return
	}

	mp.lock.Lock()
	defer mp.lock.Unlock()

	for _, subFunc := range mp.subscribers{
		//todo goroutine 종료 check
		go subFunc(message)
	}
}

